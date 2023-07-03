package main

import (
	"context"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	pb "github.com/simbayippy/orbital/kitex_gen/orbital2"
	peopleService "github.com/simbayippy/orbital/kitex_gen/orbital2/peoplesservice"

	"github.com/cloudwego/kitex/server"
)

// var _ pb.Echo = &EchoImpl{}

// EchoImpl implements the last service interface defined in the IDL.
type PeoplesServiceImpl struct{}

// Echo implements the Echo interface.
func (s *PeoplesServiceImpl) Echo(ctx context.Context, req *pb.Request) (resp *pb.Response, err error) {
	klog.Info("echo called")
	return &pb.Response{Message: req.Message}, nil
}

// EditPerson implements the PeoplesServiceImpl interface.
func (s *PeoplesServiceImpl) EditPerson(ctx context.Context, req *pb.Person) (resp *pb.Person, err error) {
	// TODO: Your code here...
	return
}

func main() {
	svr := peopleService.NewServer(new(PeoplesServiceImpl), server.WithServiceAddr(&net.TCPAddr{Port: 9000}))
	if err := svr.Run(); err != nil {
		log.Println("server stopped with error:", err)
	} else {
		log.Println("server stopped")
	}
}
