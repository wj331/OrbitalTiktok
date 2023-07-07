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

	"context"
	"encoding/json"
	"fmt"

	orbital "github.com/simbayippy/OrbitalxTiktok/RPCservers/jsonthrift/kitex_gen/orbital"

	// "github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
)

type GenericServiceImpl struct{}

// for JSON generic call

// GenericCall implements the Echo interface.
func (g *GenericServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	// use jsoniter or other json parse sdk to assert request
	m := request.(string)
	fmt.Printf("Received string:2 %s\n", m)
	//Received string:2 {"name":"","age":0}

	fmt.Print(method)
	var person orbital.Person
	if err := json.Unmarshal([]byte(m), &person); err != nil {
		klog.Fatal(err)
	}
	fmt.Print(person.Age)
	toReturn := fmt.Sprintf("{\"name\": \"%v Edited lolol\", \"age\": %d}", person.Name, person.Age)

	return toReturn, nil
}

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
