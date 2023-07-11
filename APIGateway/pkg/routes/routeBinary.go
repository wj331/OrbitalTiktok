package routes

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils"

	orbital "github.com/simbayippy/OrbitalxTiktok/APIGateway/kitex_gen/orbital"
	localUtils "github.com/simbayippy/OrbitalxTiktok/APIGateway/utils"
)

type RequestData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	rc = utils.NewThriftMessageCodec()
)

// for Binary generic call
func RegisterRouteBinaryGenericCall(h *server.Hertz) {
	h.StaticFS("/", &app.FS{Root: "./", GenerateIndexPages: true})

	v1 := h.Group("/PeopleService")
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
