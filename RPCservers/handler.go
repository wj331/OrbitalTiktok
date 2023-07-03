package main

import (
	"context"
	"encoding/json"
	"fmt"

	http "github.com/simbayippy/OrbitalxTiktok/RPCservers/kitex_gen/http"
	orbital "github.com/simbayippy/OrbitalxTiktok/RPCservers/kitex_gen/orbital"
	protopackage "github.com/simbayippy/OrbitalxTiktok/RPCservers/kitex_gen/protopackage"
	user "github.com/simbayippy/OrbitalxTiktok/RPCservers/kitex_gen/user"

	// "github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/jhump/protoreflect/desc"
)

// var _ generic.Service = &ProtoServiceImpl{}

// PeopleServiceImpl implements the last service interface defined in the IDL.
type PeopleServiceImpl struct{}
type BizServiceImpl struct{}
type UserServiceImpl struct{}
type GenericServiceImpl struct{}

type MockImpl struct{}

type ProtoServiceImpl struct {
	MessageDesc *desc.MessageDescriptor
}

// EditPerson implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl) EditPerson(ctx context.Context, person *orbital.Person) (resp *orbital.Person, err error) {
	// TODO: Your code here...
	fmt.Print("first called\n\n")
	return &orbital.Person{Name: person.Name + " edited", Age: person.Age}, nil
}

// Echo implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl) Echo(ctx context.Context, req *orbital.Request) (resp *orbital.Response, err error) {
	// TODO: Your code here...
	return &orbital.Response{Message: req.Message}, nil
}

// EditPerson2 implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl) EditPerson2(ctx context.Context, person *orbital.Person) (resp *orbital.Person, err error) {
	// TODO: Your code here...
	fmt.Print("first called haha\n\n")

	return &orbital.Person{Name: person.Name + " edited 2", Age: person.Age}, nil
}

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

// for user.thrift

// GetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	// TODO: Your code here...
	userID := req.GetUserId()

	return &user.UserResponse{
		Text:     fmt.Sprintf("User data for ID %d retrieved successfully", userID),
		HttpCode: 200,
	}, nil
}

// UpdateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	// TODO: Your code here...
	userID := req.GetUsername()

	fmt.Print(req.GetUserId())
	fmt.Print("called!\n\n")
	fmt.Print(userID)
	return &user.UserResponse{
		Text:     fmt.Sprintf("User data for %s successfully updated", userID),
		HttpCode: 200,
	}, nil
}

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

// {"name": " Edited lolol", "age": 0}
// &{response:"Received string: Msg:\"John Doe\"" 0x14000692270}
// GenericCall implements the Echo interface.
func (g *ProtoServiceImpl) GenericCall2(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	// use jsoniter or other json parse sdk to assert request
	m := request.(string)
	fmt.Printf("Received string1: %s\n", m)
	fmt.Printf("\nmethod is: %v\n", method)

	// var mockReq protopackage.MockReq
	// err = proto.Unmarshal([]byte(m), &mockReq)
	// if err != nil {
	// 	// handle error
	// 	fmt.Printf("Error unmarshalling request: %v\n", err)
	// } else {
	// 	fmt.Printf("Unmarshalled request: %+v\n", mockReq)
	// }
	// fmt.Print(mockReq)

	// // now mockReq contains the deserialized protobuf message
	// fmt.Printf("\n\nReceived MockReq: %v\n\n", mockReq)

	toReturn := &protopackage.StringResponse{ // Change here: return pointer
		Response: "hello",
	}
	return toReturn, nil
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

	toReturn := &protopackage.StringResponse{ // Change here: return pointer
		Response: "hello",
	}

	//3 jul currently error because toReturn is not a string
	return toReturn, nil
}

// ...

// func (g *ProtoServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
// 	m, ok := request.([]byte)
// 	if !ok {
// 		// handle error: request is not a byte slice
// 	}

// 	stringResponse := &protopackage.StringResponse{}

// 	// Get the descriptor from the message
// 	md := stringResponse.ProtoReflect().Descriptor()

// 	// Create a new dynamic message
// 	dm := dynamic.NewMessage(desc.ToDescriptorProto(md))

// 	// Unmarshal the request into the dynamic message
// 	if err := dm.Unmarshal(m); err != nil {
// 		klog.Fatal(err)
// 	}

// 	// Convert the dynamic message back to the original protobuf message
// 	if err := dm.ConvertTo(stringResponse); err != nil {
// 		klog.Fatal(err)
// 	}

// 	fmt.Print(stringResponse.Response)

// 	return stringResponse, nil
// }

// Test implements the MockImpl interface.
func (s *MockImpl) Test(ctx context.Context, req *protopackage.MockReq) (resp *protopackage.StringResponse, err error) {
	// TODO: Your code here...
	fmt.Print("test called\n\n")
	toAdd := req.Msg
	return &protopackage.StringResponse{Response: "response: " + toAdd}, nil
}
