package main

import (
	"context"
	"fmt"

	orbital "github.com/simbayippy/OrbitalxTiktok/RPCservers/binary/kitex_gen/orbital"
	// "github.com/cloudwego/kitex/pkg/generic"
)

// var _ generic.Service = &ProtoServiceImpl{}

// PeopleServiceImpl implements the last service interface defined in the IDL.
type PeopleServiceImpl struct{}

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
