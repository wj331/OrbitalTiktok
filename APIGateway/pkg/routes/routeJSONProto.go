package routes

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/kitex/client/genericclient"
	// imported as `cache`
	localUtils "github.com/simbayippy/OrbitalxTiktok/APIGateway/utils"
)

func RegisterRouteJSONProto(h *server.Hertz) {
	v1 := h.Group("/JSONProtoService")
	{
		v1.POST("/:version/:method", localUtils.RateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			version := c.Param("version")
			methodName := c.Param("method")

			path := string(c.Path())
			parts := strings.Split(path, "/")
			serviceName := parts[1]

			poolKey := fmt.Sprintf("%s_%s", serviceName, version)
			pool, ok := Pools[poolKey]
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
		}))
	}
}
