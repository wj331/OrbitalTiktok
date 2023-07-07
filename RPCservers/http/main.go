package main

import (
	"context"
	"log"
	"net"
	"sync"

	httpServer "github.com/simbayippy/OrbitalxTiktok/RPCservers/http/kitex_gen/http/bizservice"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/registry-nacos/registry"

	// "github.com/cloudwego/kitex/pkg/limit"

	"github.com/cloudwego/kitex/server"

	http "github.com/simbayippy/OrbitalxTiktok/RPCservers/http/kitex_gen/http"

	// "github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
)

type BizServiceImpl struct{}

// handlers
func (s *BizServiceImpl) BizMethod1(ctx context.Context, req *http.BizRequest) (resp *http.BizResponse, err error) {
	klog.Infof("BizMethod1 called, request: %#v", req)
	return &http.BizResponse{HttpCode: 200, Text: "Method1 response", Token: 1111}, nil
}

func (s *BizServiceImpl) BizMethod2(ctx context.Context, req *http.BizRequest) (resp *http.BizResponse, err error) {
	klog.Infof("BizMethod2 called, request: %#v", req)
	return &http.BizResponse{HttpCode: 200, Text: "Method2 response", Token: 2222}, nil
}

func (s *BizServiceImpl) BizMethod3(ctx context.Context, req *http.BizRequest) (resp *http.BizResponse, err error) {
	klog.Infof("BizMethod3 called, request: %#v", req)
	return &http.BizResponse{HttpCode: 200, Text: "Method3 response", Token: 3333}, nil
}

func main() {
	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// HTTP servers
	for i := 0; i < 4; i++ {
		port := 8898 + i
		svr := httpServer.NewServer(
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

	wg.Wait()
}
