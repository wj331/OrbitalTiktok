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
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	// imported as `cache`
	"github.com/patrickmn/go-cache"

	localUtils "github.com/simbayippy/OrbitalxTiktok/APIGateway/utils"
	"github.com/simbayippy/OrbitalxTiktok/APIGateway/utils/genericClients"
)

type ServiceDetails struct {
	Version  string
	FilePath string
}

// configurations
var (
	/*
		Mapping of service & IDL's
	*/
	pools = map[string]*sync.Pool{}

	services = map[string][]ServiceDetails{
		"PeopleService": {
			{Version: "v1.0.0", FilePath: "./thrift/orbital.thrift"},
			// Other versions for e.g.
			// {Version: "v2.0.0", FilePath: "./thrift/orbital_v2.thrift"},
		},
		"BizService": {
			{Version: "v1.0.0", FilePath: "./thrift/http.thrift"},
		},
		"JSONService": {
			{Version: "v1.0.0", FilePath: "./thrift/orbital.thrift"},
		},
		"JSONProtoService": {
			{Version: "v1.0.0", FilePath: "./protobuf/mock.proto"},
		},
	}

	serviceNames []string

	/*
		Caching
	*/

	// Set to true to enable caching on ALL endpoints. Default set to false for benchmarking purposes
	enableCaching = false

	// IP addresses in the cache expires after 5 minutes of no access, and the library by patrickmn automatically cleans up expired items every 6 minutes.
	limiterCache = cache.New(5*time.Minute, 6*time.Minute)

	// rate limiting numbers. set HIGH for benchmark purposes
	MaxQPS    = 10000000000000 // Each IP address how many QPS
	BurstSize = 10000000000000 // number of events that can occur at ONCE

	// cache time allowed before evicted from cache. i.e. how long stored in cache
	cacheExpiryTime = 2 * time.Second

	// cache counters
	cacheHitCount, cacheMissCount int32

	// codec to be used specifically for RegisterRouteBinaryGenericCall
	rc = utils.NewThriftMessageCodec()
)

func init() {
	// Nacos service registry
	var err error
	localUtils.NacosClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig("127.0.0.1", 8848),
			},
		},
	)
	if err != nil {
		klog.Fatalf("Failed to create Nacos client: %v", err)
	}

	// get all available service names
	serviceNames = make([]string, 0, len(services))
	for serviceName := range services {
		serviceNames = append(serviceNames, serviceName)
	}

	// Immediately adds all valid RPC instances for all services.
	localUtils.AddInitialInstance(serviceNames)

	// initializing of generic client pools upon apigateway startup
	pools = make(map[string]*sync.Pool)
	for serviceName, details := range services {
		// Set up pools for each version of this service
		for _, detail := range details {
			poolKey := serviceName + "_" + detail.Version
			pools[poolKey] = newClientPool(serviceName, detail.FilePath)
		}
	}
}

func newClientPool(serviceName string, protoFilePath string) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			var (
				cc  genericclient.Client
				err error
			)
			// TODO: change it to integer values
			switch serviceName {
			case "PeopleService":
				cc, err = genericClients.NewBinaryGenericClient(serviceName)
			case "BizService":
				cc, err = genericClients.NewHTTPGenericClient(serviceName, protoFilePath)
			case "JSONService":
				cc, err = genericClients.NewJSONGenericClient(serviceName, protoFilePath)
			case "JSONProtoService":
				cc, err = genericClients.NewJSONProtoGenericClient(serviceName, protoFilePath)
			default:
				log.Print("Invalid service name")
				return nil
			}
			if err != nil {
				log.Print("unable to generate new client")
				return nil
			}
			return cc
		},
	}
}

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	if enableCaching {
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
	}

	// register routes
	RegisterRouteJSONThrift(h)
	RegisterRouteJSONProto(h)
	RegisterRouteHTTPGenericCall(h)
	RegisterRouteBinaryGenericCall(h)
	RegisterCacheRoute(h)
	RegisterAuthRoute(h)

	// updates instances of services every minute
	localUtils.RefreshInstances(serviceNames)

	h.Spin()
}

// for JSON generic call. MAIN use case for TikTok. Basically forwards the request that the API gateway receives directly to the RPC server.
func RegisterRouteJSONThrift(h *server.Hertz) {
	v1 := h.Group("/JSONService")
	{
		v1.POST("/:version/:method", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			version := c.Param("version")
			methodName := c.Param("method")

			path := string(c.Path())
			parts := strings.Split(path, "/")
			serviceName := parts[1]

			poolKey := fmt.Sprintf("%s_%s", serviceName, version) // create the key with the service and version
			pool, ok := pools[poolKey]
			if !ok {
				c.JSON(consts.StatusBadRequest, "Invalid service name or version")
				return
			}

			// take from pool
			cc := pool.Get().(genericclient.Client)
			defer pool.Put(cc) // make sure to put the client back in the pool when done

			bodyBytes := c.GetRequest().BodyBytes()

			// checks if the request contains anything.
			// TODO: improve this
			if len(bodyBytes) == 0 {
				c.JSON(consts.StatusBadRequest, "request body is empty")
				return
			}

			jsonString := string(bodyBytes)

			resp, err := cc.GenericCall(ctx, methodName, jsonString)
			if err != nil {
				c.JSON(consts.StatusBadRequest, err)
				return
			}

			pool.Put(cc)

			c.JSON(consts.StatusOK, resp)
		}))
	}
}

func RegisterRouteJSONProto(h *server.Hertz) {
	v1 := h.Group("/JSONProtoService")
	{
		v1.POST("/:version/:method", func(ctx context.Context, c *app.RequestContext) {
			version := c.Param("version")
			methodName := c.Param("method")

			path := string(c.Path())
			parts := strings.Split(path, "/")
			serviceName := parts[1]

			poolKey := fmt.Sprintf("%s_%s", serviceName, version)
			pool, ok := pools[poolKey]
			if !ok {
				c.JSON(consts.StatusBadRequest, "Invalid service name or version")
				return
			}

			// take from pool
			cc := pool.Get().(genericclient.Client)
			defer pool.Put(cc)

			bodyBytes := c.GetRequest().BodyBytes()

			// TODO: improve this
			if len(bodyBytes) == 0 {
				c.JSON(consts.StatusBadRequest, "request body is empty")
				return
			}

			jsonString := string(bodyBytes)

			resp, err := cc.GenericCall(ctx, methodName, jsonString)
			if err != nil {
				c.JSON(consts.StatusBadRequest, err)
				return
			}

			pool.Put(cc)

			c.JSON(consts.StatusOK, resp)
		})
	}
}

// for HTTP generic call
func RegisterRouteHTTPGenericCall(h *server.Hertz) {
	v1 := h.Group("/BizService")
	{
		v1.POST("/:version/:method", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			version := c.Param("version")
			methodName := c.Param("method")

			path := string(c.Path())
			parts := strings.Split(path, "/")
			serviceName := parts[1]

			poolKey := fmt.Sprintf("%s_%s", serviceName, version) // create the key with the service and version
			pool, ok := pools[poolKey]
			if !ok {
				c.JSON(consts.StatusBadRequest, "Invalid service name or version")
				return
			}

			cc := pool.Get().(genericclient.Client)
			defer pool.Put(cc)

			data := c.GetRequest().BodyBytes()

			var jsonData map[string]interface{}
			err2 := json.Unmarshal(data, &jsonData)
			if err2 != nil {
				c.JSON(consts.StatusBadRequest, "Invalid JSON request data")
				return
			}

			url := fmt.Sprintf("http://example.com/bizservice/%s", methodName)

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
				c.JSON(consts.StatusBadRequest, err)
				return
			}
			pool.Put(cc)

			// realResp := resp.(*generic.HTTPResponse)
			c.JSON(consts.StatusOK, resp)
		}))
	}
}

// for Binary generic call
func RegisterRouteBinaryGenericCall(h *server.Hertz) {
	h.StaticFS("/", &app.FS{Root: "./", GenerateIndexPages: true})

	type RequestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	v1 := h.Group("/PeopleService")
	{
		v1.POST("/:version/:method", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			version := c.Param("version")
			methodName := c.Param("method")

			path := string(c.Path())
			parts := strings.Split(path, "/")
			serviceName := parts[1]

			poolKey := fmt.Sprintf("%s_%s", serviceName, version) // create the key with the service and version
			pool, ok := pools[poolKey]
			if !ok {
				c.JSON(consts.StatusBadRequest, "Invalid service name or version")
				return
			}

			cc := pool.Get().(genericclient.Client)
			defer pool.Put(cc) // make sure to put the client back in the pool when done

			var requestData RequestData

			// Bind and parse the request body into the requestData struct
			err := c.Bind(&requestData)
			if err != nil {
				// Handle error if the request body parsing fails
				c.JSON(consts.StatusBadRequest, "Invalid request body")
				return
			}

			// Access the data from the request body
			name := requestData.Name
			age := requestData.Age

			// Create a Person object from the request data
			person := &orbital.PeopleServiceEditPersonArgs{Person: &orbital.Person{Name: name, Age: int32(age)}}

			// method "editPerson" has to follow exactly the same spelling (capital) as in .thrift file
			buf, err := rc.Encode(methodName, thrift.CALL, 0, person)
			if err != nil {
				klog.Errorf("failed to encode: %w", err)
				return
			}

			// GenericCall feature of kitex to "generically" call this method in the RPC server
			resp, err := cc.GenericCall(ctx, methodName, buf)
			if err != nil {
				c.JSON(consts.StatusBadRequest, err)
				return
			}

			pool.Put(cc)

			// **DESERIALIZING**

			// creating an empty result struct for the result to be decoded into
			result := &orbital.PeopleServiceEditPersonResult{}
			_, _, err = rc.Decode(resp.([]byte), result)
			if err != nil {
				klog.Fatal(err)
				return
			}

			// Send a response
			c.JSON(consts.StatusOK, "all okay")
		}))
	}
}

// RegisterCacheRoute: to demonstrate caching
// Caching: as long as request URI is the same -> immediately returns cached response -> improve response times
func RegisterCacheRoute(h *server.Hertz) {
	v1 := h.Group("/cache")
	{
		// in apache benchmarking -> result cached -> used for rest of benchmark. i.e. only one full request made, rest of
		// requests are done using the first requests results
		v1.POST("/post", rateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			response := "caching....."

			c.String(consts.StatusOK, response)
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
