package main

import (
	"context"
	"encoding/json"
	"fmt"

	protopackage "github.com/simbayippy/OrbitalxTiktok/RPCservers/kitex_gen/protopackage"

	// "github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/jhump/protoreflect/desc"
)

type ProtoServiceImpl struct {
	MessageDesc *desc.MessageDescriptor
}

// request is a string representation of the protobuf message
func (g *ProtoServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	// use jsoniter or other json parse sdk to assert request
	// use jsoniter or other json parse sdk to assert request
	m := request.(string)
	fmt.Printf("Received string:2 %s\n", m)
	//Received string:2 {"name":"","age":0}
	var mockReq protopackage.MockReq
	if err := json.Unmarshal([]byte(m), &mockReq); err != nil {
		klog.Fatal(err)
	}
	fmt.Print(mockReq.Msg)
	fmt.Print(mockReq.StrList)

	toReturn := &protopackage.StringResponse{
		Response: "hello",
	}

	return toReturn, nil
}
