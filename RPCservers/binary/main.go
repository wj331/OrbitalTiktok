package main

import (
	"context"
	"log"
	"net"
	"sync"

	orbitalServer "github.com/simbayippy/OrbitalxTiktok/RPCservers/binary/kitex_gen/orbital/peopleservice"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"

	// "github.com/cloudwego/kitex/pkg/limit"

	orbital "github.com/simbayippy/OrbitalxTiktok/RPCservers/binary/kitex_gen/orbital"
)

type PeopleServiceImpl struct{}

// EditPerson implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl) EditPerson(ctx context.Context, person *orbital.Person) (resp *orbital.Person, err error) {
	return &orbital.Person{Name: person.Name + " edited", Age: person.Age}, nil
}

// Echo implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl) Echo(ctx context.Context, req *orbital.Request) (resp *orbital.Response, err error) {
	return &orbital.Response{Message: req.Message}, nil
}

// EditPerson2 implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl) EditPerson2(ctx context.Context, person *orbital.Person) (resp *orbital.Person, err error) {
	return &orbital.Person{Name: person.Name + " edited 2", Age: person.Age}, nil
}

func main() {
	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// binary servers
	for i := 0; i < 4; i++ {
		port := 8888 + i
		svr := orbitalServer.NewServer(
			new(PeopleServiceImpl),
			server.WithServiceAddr(&net.TCPAddr{Port: port}),
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "PeopleService"}),
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
