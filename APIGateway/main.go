package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/basic_auth"
	"github.com/cloudwego/hertz/pkg/app/server"

	// hertzUtils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	hertzCache "github.com/hertz-contrib/cache"
	"github.com/hertz-contrib/cache/persist"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils"

	orbital "github.com/simbayippy/OrbitalxTiktok/APIGateway/kitex_gen/orbital"

	"github.com/nacos-group/nacos-sdk-go/clients"
	// "github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	// imported as `cache`
	"github.com/patrickmn/go-cache"

	// new
	localUtils "github.com/simbayippy/OrbitalxTiktok/APIGateway/utils"
	"github.com/simbayippy/OrbitalxTiktok/APIGateway/utils/genericClients"
)

// local variables for **OPTIMIZATION**
var (
	// creating a pool of generic binary clients
	ccBinaryPool = &sync.Pool{
		New: func() interface{} {
			// TODO, create a new pool of clients that serve other services too
			cc, err := genericClients.NewBinaryGenericClient("PeopleService")
			if err != nil {
				log.Print("unable to generate new client")
				return nil
			}
			return cc
		},
	}

	// creating a pool of generic HTTP clients for service: UserService
	ccHTTPUserServicePool = &sync.Pool{
		New: func() interface{} {
			// since requires a NewThriftFileProvider(file_path)
			cc, err := genericClients.NewHTTPGenericClient("UserService", "./thrift/user.thrift")
			if err != nil {
				log.Print("unable to generate new client")
				return nil
			}
			return cc
		},
	}

	// creating a pool of generic HTTP clients for service: UserService
	ccHTTPBizServicePool = &sync.Pool{
		New: func() interface{} {
			// since requires a NewThriftFileProvider(file_path)
			cc, err := genericClients.NewHTTPGenericClient("BizService", "./thrift/http.thrift")
			if err != nil {
				log.Print("unable to generate new client")
				return nil
			}
			return cc
		},
	}

	// creating a pool of generic HTTP clients for service: UserService
	ccJSONPool = &sync.Pool{
		New: func() interface{} {
			// due to similarities of how servers are implemented, reusing orbital.thrift file
			cc, err := genericClients.NewJSONGenericClient("JSONService", "./thrift/orbital.thrift")
			if err != nil {
				log.Print("unable to generate new client")
				return nil
			}
			return cc
		},
	}

	ccJSONProtoPool = &sync.Pool{
		New: func() interface{} {
			// due to similarities of how servers are implemented, reusing orbital.thrift file
			cc, err := genericClients.NewJSONProtoGenericClient("JSONProtoService", "./protobuf/mock.proto")
			if err != nil {
				log.Print("unable to generate new client")
				return nil
			}
			return cc
		},
	}

	// then a mapping of all the pools. O(1) to get the pool to use in GO, so it's very efficient.
	// small capital letter for path of apigateway such as api/peopleservice/methodname
	pools = map[string]*sync.Pool{
		"peopleservice":    ccBinaryPool,
		"userservice":      ccHTTPUserServicePool,
		"bizservice":       ccHTTPBizServicePool,
		"jsonservice":      ccJSONPool,
		"jsonprotoservice": ccJSONProtoPool,
	}

	// local codec to be used for RegisterRouteBinaryGenericCall/ caching
	rc = utils.NewThriftMessageCodec()

	// IP addresses in the cache expires after 5 minutes of no access, and the library by patrickmn automatically cleans up expired items every 6 minutes.
	limiterCache = cache.New(5*time.Minute, 6*time.Minute)

	// TODO: rate limiting numbers
	MaxQPS    = 1000000000 // Each IP address how many QPS
	BurstSize = 1000000000 // number of events that can occur at ONCE. set HIGH for benchmark purposes

	// TODO: cache time allowed before evicted from cache. i.e. how long stored in cache
	cacheExpiryTime = 2 * time.Second

	// cache counters
	cacheHitCount, cacheMissCount int32

	services = []string{"PeopleService", "UserService", "BizService", "JSONService", "JSONProtoService"}
)

func init() {

	// initializes the nacos client: used for service registry (RPC servers) and discovery (this API gateway)
	var err error
	// New
	localUtils.NacosClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig("127.0.0.1", 8848),
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Nacos client: %v", err)
	}

	//new
	localUtils.AddInitialInstance(services)
}

func main() {
	// using Hertz library to generate new Hertz server
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	/*
		Uncomment the following lines if want to use caching on ALL endpoints.
		Commented out to see performance in AB, as with caching technically only 1 request will be sent to RPC.
	*/

	// memoryStore := persist.NewMemoryStore(1 * time.Minute)
	// h.Use(hertzCache.NewCacheByRequestURI(
	// 	memoryStore,
	// 	2*time.Second,
	// 	hertzCache.WithOnHitCache(func(ctx context.Context, c *app.RequestContext) {
	// 		atomic.AddInt32(&cacheHitCount, 1)
	// 	}),
	// 	hertzCache.WithOnMissCache(func(ctx context.Context, c *app.RequestContext) {
	// 		atomic.AddInt32(&cacheMissCount, 1)
	// 	}),
	// ))

	// register routes

	RegisterRouteJSONGenericCall(h)
	RegisterRouteJSONProto(h)
	RegisterRouteBinaryGenericCall(h)
	RegisterRouteHTTPGenericCall(h)
	RegisterCacheRoute(h)
	RegisterAuthRoute(h)

	// new
	localUtils.RefreshInstances(services)

	h.Spin()
}

// IMPORTANT (ROUTING PRIORITY): static route > parametric route > wildcard route

// for JSON generic call. MAIN use case for TikTok. Basically forwards the request that the API gateway receives directly to the RPC server.
func RegisterRouteJSONGenericCall(h *server.Hertz) {
	v1 := h.Group("/jsonservice")
	{
		v1.POST("/:method", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			methodName := c.Param("method")

			path := string(c.Path())
			parts := strings.Split(path, "/")
			// which here is userservice
			serviceName := parts[1]

			pool, ok := pools[serviceName]
			if !ok {
				// handle the case where there is no pool for the given service name
				c.String(consts.StatusBadRequest, "Invalid service name")
				return
			}

			// take from pool
			cc := pool.Get().(genericclient.Client)
			defer pool.Put(cc) // make sure to put the client back in the pool when done

			bodyBytes := c.GetRequest().BodyBytes()

			// checks if the request contains anything.
			// TODO: improve this
			if len(bodyBytes) == 0 {
				c.String(consts.StatusBadRequest, "request body is empty")
				return
			}

			jsonString := string(bodyBytes)

			// methodName string is case sensitive as per implementation of RPC server
			resp, err := cc.GenericCall(ctx, methodName, jsonString)
			if err != nil {
				// if method name does not exist, logged here.
				klog.Errorf("generic call failed: %v", err)
				return
			}

			// putting client back into pool after generic call is called
			pool.Put(cc)

			c.String(consts.StatusOK, "response is: %v", resp)
		}))
	}
}

func RegisterRouteJSONProto(h *server.Hertz) {
	v1 := h.Group("/jsonprotoservice")
	{
		v1.POST("/:method", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			methodName := c.Param("method")

			path := string(c.Path())
			parts := strings.Split(path, "/")
			// which here is userservice
			serviceName := parts[1]

			pool, ok := pools[serviceName]
			if !ok {
				// handle the case where there is no pool for the given service name
				c.String(consts.StatusBadRequest, "Invalid service name")
				return
			}

			// take from pool
			cc := pool.Get().(genericclient.Client)
			defer pool.Put(cc) // make sure to put the client back in the pool when done

			// request should have a json string in it
			bodyBytes := c.GetRequest().BodyBytes()

			// convert bodyBytes to string, should be in json format
			jsonString := string(bodyBytes)

			// fmt.Printf("\nLog: (API Gateway) data received in string form is: %v\n", jsonString)

			// methodName string is case sensitive as per implementation of RPC server
			resp, err := cc.GenericCall(ctx, methodName, jsonString)
			if err != nil {
				// if method name does not exist, logged here.
				klog.Errorf("generic call failed: %v", err)
				return
			}

			// putting client back into pool after generic call is called
			pool.Put(cc)

			c.String(consts.StatusOK, "response is: %v", resp)
		}))
	}
}

type RequestData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// for Binary generic call
func RegisterRouteBinaryGenericCall(h *server.Hertz) {
	h.StaticFS("/", &app.FS{Root: "./", GenerateIndexPages: true})

	h.GET("/get", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "get")
	}))

	h.POST("/post", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
		var requestData RequestData

		// Bind and parse the request body into the requestData struct
		err := c.Bind(&requestData)
		if err != nil {
			// Handle error if the request body parsing fails
			c.String(consts.StatusBadRequest, "Invalid request body")
			return
		}

		// Access the data from the request body
		name := requestData.Name
		age := requestData.Age

		// Create a Person object from the request data
		person := &orbital.PeopleServiceEditPersonArgs{Person: &orbital.Person{Name: name, Age: int32(age)}}

		// fmt.Printf("Name is: %v, Age is %v\n\n", name, age)

		// **SERIALIZING**

		// **OPTIMIZATION** created a pool of clients and getting from the pool instead of generating new client every request
		cc := ccBinaryPool.Get().(genericclient.Client)
		defer ccBinaryPool.Put(cc) // make sure to put the client back in the pool when done

		// method "editPerson" has to follow exactly the same spelling (capital) as in .thrift file
		buf, err := rc.Encode("editPerson", thrift.CALL, 0, person)
		if err != nil {
			klog.Errorf("failed to encode: %w", err)
			return
		}

		// GenericCall feature of kitex to "generically" call this method in the RPC server
		resp, err := cc.GenericCall(ctx, "editPerson", buf)
		if err != nil {
			klog.Errorf("method call editPerson failed: %w", err)
			return
		}

		ccBinaryPool.Put(cc)

		// **DESERIALIZING**

		// creating an empty result struct for the result to be decoded into
		result := &orbital.PeopleServiceEditPersonResult{}
		_, _, err = rc.Decode(resp.([]byte), result)
		if err != nil {
			klog.Fatal(err)
			return
		}

		// Send a response
		c.String(consts.StatusOK, "all okay")
	}))
}

// for HTTP generic call. also dependent on how .thrift file is implemented!!
func RegisterRouteHTTPGenericCall(h *server.Hertz) {
	v2 := h.Group("/bizservice")
	{
		v2.POST("/:method", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			method := c.Param("method")

			path := string(c.Path())
			parts := strings.Split(path, "/")
			// which here is bizservice
			serviceName := parts[1]

			pool, ok := pools[serviceName]
			if !ok {
				// handle the case where there is no pool for the given service name
				c.String(consts.StatusBadRequest, "Invalid service name")
				return
			}

			cc := pool.Get().(genericclient.Client)
			defer pool.Put(cc) // make sure to put the client back in the pool when done

			data := c.GetRequest().BodyBytes()

			var jsonData map[string]interface{}
			err2 := json.Unmarshal(data, &jsonData)
			if err2 != nil {
				c.String(consts.StatusBadRequest, "Invalid JSON request data")
				return
			}

			url := fmt.Sprintf("http://example.com/bizservice/%s", method)

			// cannot get the current http request from argument c, and not a drastic overhead to just create a new http request.
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			if err != nil {
				klog.Errorf("new http request failed: %v", err)
				return
			}
			req.Header.Set("token", "3")
			customReq, err := generic.FromHTTPRequest(req)
			if err != nil {
				klog.Errorf("convert request failed: %v", err)
				return
			}
			resp, err := cc.GenericCall(ctx, "", customReq)
			if err != nil {
				klog.Errorf("generic call failed: %v", err)
				return
			}
			pool.Put(cc)

			realResp := resp.(*generic.HTTPResponse)
			c.String(consts.StatusOK, "UpdateUser response, status code: %v, body: %v\n", realResp.StatusCode, realResp.Body)
		}))
	}
}

// RegisterCacheRoute: to demonstrate caching
// Caching: as long as request URI is the same -> immediately returns cached response -> improve response times
func RegisterCacheRoute(h *server.Hertz) {
	memoryStore := persist.NewMemoryStore(1 * time.Minute)
	h.Use(hertzCache.NewCacheByRequestURI(
		memoryStore,
		cacheExpiryTime,
		hertzCache.WithOnHitCache(func(ctx context.Context, c *app.RequestContext) {
			atomic.AddInt32(&cacheHitCount, 1)
		}),
		hertzCache.WithOnMissCache(func(ctx context.Context, c *app.RequestContext) {
			atomic.AddInt32(&cacheMissCount, 1)
		}),
	))
	v1 := h.Group("/cache")
	{
		// exact same method as /post above. just to demonstrate power of caching
		// in apache bench -> request if forwarded to backend RPC only ONCE -> result cached -> used for rest of benchmark.
		v1.POST("/post", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			var requestData RequestData

			err := c.Bind(&requestData)
			if err != nil {
				c.String(consts.StatusBadRequest, "Invalid request body")
				return
			}

			name := requestData.Name
			age := requestData.Age

			person := &orbital.PeopleServiceEditPersonArgs{Person: &orbital.Person{Name: name, Age: int32(age)}}

			cc := ccBinaryPool.Get().(genericclient.Client)
			defer ccBinaryPool.Put(cc)

			buf, err := rc.Encode("editPerson", thrift.CALL, 0, person)
			if err != nil {
				return
			}

			resp, err := cc.GenericCall(ctx, "editPerson", buf)
			if err != nil {
				klog.Errorf("method call editPerson failed: %w", err)
				return
			}

			ccBinaryPool.Put(cc)

			result := &orbital.PeopleServiceEditPersonResult{}
			_, _, err = rc.Decode(resp.([]byte), result)
			if err != nil {
				log.Fatal(err)
				return
			}

			c.String(consts.StatusOK, "all okay")
		}))
		// if response is in cache
		v1.GET("/get_hit_count", func(ctx context.Context, c *app.RequestContext) {
			c.String(200, fmt.Sprintf("total hit count: %d", cacheHitCount))
		})
		v1.GET("/get_miss_count", func(ctx context.Context, c *app.RequestContext) {
			c.String(200, fmt.Sprintf("total miss count: %d", cacheMissCount))
		})
	}
}

// RegisterAuthRoute group route
func RegisterAuthRoute(h *server.Hertz) {
	h.Use(basic_auth.BasicAuth(map[string]string{
		"username1": "password1",
		"username2": "password2",
	}))
	v1 := h.Group("/auth")
	{
		// "get" is only returned if username and password is included in Header of client request
		v1.GET("/get", func(ctx context.Context, c *app.RequestContext) {
			c.String(consts.StatusOK, "get")
		})
	}
}

/*
	Helper functions: MOVED TO THEIR OWN FILES: discovery, genericClient and ratelimt. Remain in same package
*/
