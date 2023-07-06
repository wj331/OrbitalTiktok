package main

import (
	"log"
	"net"

	"sync"

	orbital "github.com/simbayippy/OrbitalxTiktok/RPCservers/binary/kitex_gen/orbital/peopleservice"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/registry-nacos/registry"

	// "github.com/cloudwego/kitex/pkg/limit"

	"github.com/cloudwego/kitex/server"
)

func main() {

	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// binary servers
	// Create multiple servers to demonstrate load balancing
	for i := 0; i < 1; i++ {
		port := 8888 + i
		svr := orbital.NewServer(
			new(PeopleServiceImpl),
			server.WithServiceAddr(&net.TCPAddr{Port: port}),
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "PeopleService"}),
			server.WithRegistry(r),
			// server.WithLimit(&limit.Option{MaxConnections: 1000000, MaxQPS: 100000}),
		)

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
