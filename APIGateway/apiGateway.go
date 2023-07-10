package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/basic_auth"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/hertz-contrib/cache/persist"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	// imported as `cache`

	"github.com/simbayippy/OrbitalxTiktok/APIGateway/pkg/genericClients"
	"github.com/simbayippy/OrbitalxTiktok/APIGateway/pkg/routes"
	localUtils "github.com/simbayippy/OrbitalxTiktok/APIGateway/utils"
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

	// serviceNames data will be filled upon initialization
	serviceNames []string

	/*
		Caching
	*/
	enableCaching = false // Set true to enable caching on ALL endpoints. Default set to false for benchmarking purposes

	// cache time allowed before evicted from cache. i.e. how long stored in cache
	cacheExpiryTime = 2 * time.Second
	memoryStore     = persist.NewMemoryStore(1 * time.Minute)

	/*
		Rate limiting
	*/
	MaxQPS    = 10000000000000 // default set to 10,000
	BurstSize = 10000000000000 // default set to 100
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
	localUtils.SetRatelimits(MaxQPS, BurstSize)

	// Initializes the generic client pool
	routes.Pools = genericClients.InitGenericClientPool(services)
}

func main() {
	h := server.Default(server.WithHostPorts(APIGatewayHostPort))

	if enableCaching {
		cacheDetails := localUtils.CachingDetails(memoryStore, cacheExpiryTime)
		h.Use(cacheDetails)
	}

	// register routes
	routes.RegisterRouteJSONThrift(h)
	routes.RegisterRouteJSONProto(h)
	routes.RegisterRouteHTTPGenericCall(h)
	routes.RegisterRouteBinaryGenericCall(h)
	RegisterCacheRoute(h)
	RegisterAuthRoute(h)

	// updates instances of services every minute
	localUtils.RefreshInstances(serviceNames)

	h.Spin()
}

/*
	Left the basic route handlers as per documentation of Kitex here.
	Abstracted away main implementations to pkg/routes
*/

// RegisterCacheRoute: to demonstrate caching
// Caching: as long as request URI is the same -> immediately returns cached response -> improve response times
func RegisterCacheRoute(h *server.Hertz) {
	v1 := h.Group("/cache")
	{
		// in apache benchmarking -> result cached -> used for rest of benchmark. i.e. only one full request made, rest of
		// requests are done using the first requests results
		v1.POST("/post", func(ctx context.Context, c *app.RequestContext) {
			response := "caching....."

			c.String(consts.StatusOK, response)
		})
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
