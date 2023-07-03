// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"log"

	"orbital/kitex_gen/orbital"
	"orbital/kitex_gen/orbital/peopleservice"

	"github.com/cloudwego/kitex/client"
)

func main() {
	c, err := peopleservice.NewClient("hello", client.WithHostPorts("0.0.0.0:8889"))
	if err != nil {
		log.Fatal(err)
	}
	req := &orbital.Request{Message: "my request"}
	resp, err := c.Echo(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	req1 := &orbital.Person{Name: "jack", Age: int32(21)}
	resp1, err := c.EditPerson(context.Background(), req1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp1)
}
