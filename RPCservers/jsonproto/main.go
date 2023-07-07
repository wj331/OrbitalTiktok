package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"sync"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server/genericserver"
	"github.com/kitex-contrib/registry-nacos/registry"

	// "github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"

	protopackage "github.com/simbayippy/OrbitalxTiktok/RPCservers/kitex_gen/protopackage"

	"github.com/jhump/protoreflect/desc"
)

type ProtoServiceImpl struct {
	MessageDesc *desc.MessageDescriptor
}

// request is a string representation of the protobuf message
func (g *ProtoServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	m := request.(string)
	var mockReq protopackage.MockReq
	if err := json.Unmarshal([]byte(m), &mockReq); err != nil {
		klog.Fatal(err)
	}

	// send back a simple response
	toReturn := &protopackage.StringResponse{
		Response: "hello",
	}

	return toReturn, nil
}

func main() {
	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	protoFilePath := "./proto/mock.proto"

	p, err := generic.NewPbFileProvider(protoFilePath)
	if err != nil {
		panic(err)
	}

	g, err := generic.JSONProtoGeneric(p)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 4; i++ {
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
