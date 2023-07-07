package RPCservers

import (
	"log"
	"net"

	"sync"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server/genericserver"
	"github.com/kitex-contrib/registry-nacos/registry"

	// "github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/server"
)

func main() {

	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// jsonthrift server
	p, err := generic.NewThriftFileProvider("./thrift/orbital.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1; i++ {
		port := 8918 + i
		svr := genericserver.NewServer(
			new(GenericServiceImpl),
			g,
			server.WithServiceAddr(&net.TCPAddr{Port: port}),
			server.WithRegistry(r),
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "JSONService"}))

		wg.Add(1)

		go func() {
			defer wg.Done()
			if err := svr.Run(); err != nil {
				log.Printf("server at port %d stopped with error: %v\n", port, err)
			} else {
				log.Printf("server at port %d stopped\n", port)
			}
		}()
	}

	wg.Wait()
}
