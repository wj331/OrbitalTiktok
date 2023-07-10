package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/basic_auth"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/cloudwego/hertz/pkg/protocol/consts"

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

	localUtils "github.com/simbayippy/OrbitalxTiktok/APIGateway/utils"
	"github.com/simbayippy/OrbitalxTiktok/APIGateway/utils/genericClients"
)

// configurations
var (
	/*
		Ports & Addresses
	*/
	APIGatewayHostPort = "127.0.0.1:8080" // where this api gateway will be hosted
	nacosIpAddr        = "127.0.0.1"
	nacosPortNum       = 8848

	/*
		Mapping of service & IDL's
	*/
	services = map[string][]genericClients.ServiceDetails{
		"PeopleService": {
			{Version: "v1.0.0", FilePath: "./thrift/orbital.thrift", GenericClientType: 1},
			// Other versions for e.g.
			// {Version: "v2.0.0", FilePath: "./thrift/orbital_v2.thrift"},
		},
		"BizService": {
			{Version: "v1.0.0", FilePath: "./thrift/http.thrift", GenericClientType: 2},
		},
		"JSONService": {
			{Version: "v1.0.0", FilePath: "./thrift/orbital.thrift", GenericClientType: 3},
		},
		"JSONProtoService": {
			{Version: "v1.0.0", FilePath: "./protobuf/mock.proto", GenericClientType: 4},
		},
	}

	// Following will be filled upon initialization
	pools        = map[string]*sync.Pool{}
	serviceNames []string

	/*
		Caching
	*/
	// Set true to enable caching on ALL endpoints. Default set to false for benchmarking purposes
	enableCaching = false

	// cache time allowed before evicted from cache. i.e. how long stored in cache
	cacheExpiryTime = 2 * time.Second
	memoryStore     = persist.NewMemoryStore(1 * time.Minute)

	/*
		Rate limiting
	*/
	MaxQPS    = 10000000000000
	BurstSize = 10000000000000
)

func init() {
	// Nacos service registry
	var err error
	localUtils.NacosClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig(nacosIpAddr, uint64(nacosPortNum)),
			},
		},
	)
	if err != nil {
		klog.Fatalf("Failed to create Nacos client: %v", err)
	}

	// get all available services as per configuration
	serviceNames = make([]string, 0, len(services))
	for serviceName := range services {
		serviceNames = append(serviceNames, serviceName)
	}

	// Initializes/adds all valid RPC instances for all services.
	localUtils.AddInitialInstance(serviceNames)

	// Initializes the generic client pool
	pools = genericClients.InitGenericClientPool(services)
}

func main() {
	h := server.Default(server.WithHostPorts(APIGatewayHostPort))

	if enableCaching {
		cacheDetails := localUtils.CachingDetails(memoryStore, cacheExpiryTime)
		h.Use(cacheDetails)
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
	rc := utils.NewThriftMessageCodec()

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
			c.String(200, fmt.Sprintf("total hit count: %d", localUtils.CacheHitCount))
		})
		v1.GET("/get_miss_count", func(ctx context.Context, c *app.RequestContext) {
			c.String(200, fmt.Sprintf("total miss count: %d", localUtils.CacheMissCount))
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
