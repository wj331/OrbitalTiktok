package main

import (
	"context"

	http "github.com/simbayippy/OrbitalxTiktok/RPCservers/http/kitex_gen/http"

	// "github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
)

// var _ generic.Service = &ProtoServiceImpl{}

// PeopleServiceImpl implements the last service interface defined in the IDL.
type BizServiceImpl struct{}

// handlers for http.thrift

// BizMethod1 implements the BizServiceImpl interface.
func (s *BizServiceImpl) BizMethod1(ctx context.Context, req *http.BizRequest) (resp *http.BizResponse, err error) {
	klog.Infof("BizMethod1 called, request: %#v", req)
	return &http.BizResponse{HttpCode: 200, Text: "Method1 response", Token: 1111}, nil
}

// BizMethod2 implements the BizServiceImpl interface.
func (s *BizServiceImpl) BizMethod2(ctx context.Context, req *http.BizRequest) (resp *http.BizResponse, err error) {
	klog.Infof("BizMethod2 called, request: %#v", req)
	return &http.BizResponse{HttpCode: 200, Text: "Method2 response", Token: 2222}, nil
}

// BizMethod3 implements the BizServiceImpl interface.
func (s *BizServiceImpl) BizMethod3(ctx context.Context, req *http.BizRequest) (resp *http.BizResponse, err error) {
	klog.Infof("BizMethod3 called, request: %#v", req)
	return &http.BizResponse{HttpCode: 200, Text: "Method3 response", Token: 3333}, nil
}
