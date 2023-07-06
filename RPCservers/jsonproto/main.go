package main

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

	// jsonproto server
	protoFilePath := "./proto/mock.proto"

	p, err := generic.NewPbFileProvider(protoFilePath)
	if err != nil {
		panic(err)
	}

	g, err3 := generic.JSONProtoGeneric(p)
	if err3 != nil {
		panic(err3)
	}

	for i := 0; i < 1; i++ {
		port := 8928 + i
		svr := genericserver.NewServer(
			new(ProtoServiceImpl),
			g,
			server.WithServiceAddr(&net.TCPAddr{Port: port}),
			server.WithRegistry(r),
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "JSONProtoService"}))

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
