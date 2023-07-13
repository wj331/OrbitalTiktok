package main

import (
	"context"
	// "encoding/json"
	"fmt"
	// "os"
	// "os/signal"
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

	// "github.com/fsnotify/fsnotify"

	"github.com/simbayippy/OrbitalxTiktok/APIGateway/pkg/routes"

	"github.com/simbayippy/OrbitalxTiktok/APIGateway/pkg/serviceDetails"
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

	// 1) Initializes the generic client pool, using the service configuration file
	serviceDetails.InitServiceMapping()

	// 2) Get the string list of services available as per configuration
	serviceNames := serviceDetails.GetServiceNames()

	// 3) Get valid RPC instances for all services.
	localUtils.AddInitialInstance(serviceNames)
	localUtils.SetRatelimits(MaxQPS, BurstSize)
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
	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				// gets the services available / checks if there is a change in services
				serviceNames := serviceDetails.GetServiceNames()

				// refreshes & updates list of RPC instances, even if there is no change in services
				// -> even if no change in services, available RPC instances may change
				localUtils.RefreshInstances(serviceNames)
			}
		}
	}()

	// continuously check for changes in config file
	go serviceDetails.WatchServiceChanges()

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
