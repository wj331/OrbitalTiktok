package main

import (
	"log"
	"net"

	"sync"

	orbital "orbital/kitex_gen/orbital/peopleservice"

	http "orbital/kitex_gen/http/bizservice"

	user "orbital/kitex_gen/user/userservice"

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

	// HTTP server
	for i := 0; i < 1; i++ {
		port := 8898 + i
		svr := http.NewServer(
			new(BizServiceImpl),
			server.WithServiceAddr(&net.TCPAddr{Port: port}),
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "BizService"}),
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

	for i := 0; i < 1; i++ {
		port := 8908 + i
		svr := user.NewServer(
			new(UserServiceImpl),
			server.WithServiceAddr(&net.TCPAddr{Port: port}),
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "UserService"}),
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

	// jsonproto server
	protoFilePath := "./proto/mock.proto"

	p2, err2 := generic.NewPbFileProvider(protoFilePath)
	if err2 != nil {
		panic(err)
	}

	g, err3 := generic.JSONProtoGeneric(p2)
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
