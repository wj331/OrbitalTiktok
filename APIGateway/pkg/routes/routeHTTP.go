package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"

	localUtils "github.com/simbayippy/OrbitalxTiktok/APIGateway/utils"
)

// for HTTP generic call
func RegisterRouteHTTPGenericCall(h *server.Hertz) {
	v1 := h.Group("/BizService")
	{
		v1.POST("/:version/:method", localUtils.RateLimitMiddleware(func(ctx context.Context, c *app.RequestContext) {
			version := c.Param("version")
			methodName := c.Param("method")

			path := string(c.Path())
			parts := strings.Split(path, "/")
			serviceName := parts[1]

			poolKey := fmt.Sprintf("%s_%s", serviceName, version) // create the key with the service and version
			pool, ok := Pools[poolKey]
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
