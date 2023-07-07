package apiGateway

import (
	"context"
	"fmt"

	orbital "github.com/simbayippy/OrbitalxTiktok/APIGateway/kitex_gen/orbital"
	protopackage "github.com/simbayippy/OrbitalxTiktok/APIGateway/kitex_gen/protopackage"
)

var (
	_ orbital.PeopleService = &PeopleServiceImpl{}
	_ orbital.PeopleService = &PeopleServiceImpl1{}
)

// PeopleServiceImpl implements the last service interface defined in the IDL.
type PeopleServiceImpl struct{}
type PeopleServiceImpl1 struct{}
type PeoplesServiceImpl struct{}

type BizServiceImpl struct{}

type MockImpl struct{}

// EditPerson implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl) EditPerson(ctx context.Context, person *orbital.Person) (resp *orbital.Person, err error) {
	// TODO: Your code here...
	fmt.Print("first called\n\n")

	return person, nil
}

// EditPerson implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl1) EditPerson(ctx context.Context, person *orbital.Person) (resp *orbital.Person, err error) {
	// TODO: Your code here...
	fmt.Print("second called\n\n")

	return person, nil
}

// Echo implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl) Echo(ctx context.Context, req *orbital.Request) (resp *orbital.Response, err error) {
	// TODO: Your code here...
	return &orbital.Response{Message: req.Message}, nil
}

// Echo implements the PeopleServiceImpl interface.
func (s *PeopleServiceImpl1) Echo(ctx context.Context, req *orbital.Request) (resp *orbital.Response, err error) {
	// TODO: Your code here...
	return &orbital.Response{Message: req.Message}, nil
}

// Test implements the MockImpl interface.
func (s *MockImpl) Test(ctx context.Context, req *protopackage.MockReq) (resp *protopackage.StringResponse, err error) {
	// TODO: Your code here...
	return
}
