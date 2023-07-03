package main

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
	pb "github.com/simbayippy/OrbitalxTiktok/APIGateway/kitex_gen/orbital2"
	peopleService "github.com/simbayippy/OrbitalxTiktok/APIGateway/kitex_gen/orbital2/peoplesservice"
)

func main() {
	client, err := peopleService.NewClient("peoplesservice", client.WithHostPorts("0.0.0.0:9000"))
	if err != nil {
		log.Fatal(err)
	}
	for {
		req := &pb.Request{Message: "my request"}
		resp, err := client.Echo(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
		time.Sleep(time.Second)
	}
}
