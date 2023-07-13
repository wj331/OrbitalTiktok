package main

import (
	"bytes"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pb "github.com/simbayippy/OrbitalxTiktok/APIGateway/kitex_gen/protopackage"
	"google.golang.org/protobuf/encoding/protojson"
)

func sendProtoRequest() {
	// Prepare your request
	req := &pb.MockReq{
		Msg: "MockReq",
		StrMap: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
		StrList: []string{"str1", "str2"},
	}

	// Marshal your request into binary data
	// data, err := proto.Marshal(req)
	// if err != nil {
	// 	log.Fatalf("Failed to encode: %v", err)
	// }
	m := protojson.MarshalOptions{}
	data, err := m.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}

	// Prepare and do the HTTP request
	httpReq, err := http.NewRequest("POST", "http://127.0.0.1:8080/jsonprotoservice/Test", bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	httpReq.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Fatalf("Failed to do request: %v", err)
	}
	defer resp.Body.Close()

	// Read and unmarshal the response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	bodyString := string(bodyBytes)

	// Now respMsg contains your response
	fmt.Printf("%v", bodyString)
}

func main() {
	sendProtoRequest()
}
