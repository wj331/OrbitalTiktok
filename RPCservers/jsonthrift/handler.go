package main

import (
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
