package main

import (
	"errors"

	"google.golang.org/protobuf/proto"

	pb "github.com/simbayippy/orbital/kitex_gen/orbital2"
)

type ProtobufMethodCall struct {
	Method string
	Msg    proto.Message
}

type ProtobufMessageCodec struct{}

func NewProtobufMessageCodec() *ProtobufMessageCodec {
	return &ProtobufMessageCodec{}
}

// func (p *ProtobufMessageCodec) Encode(method string, msg proto.Message) ([]byte, error) {
// 	if method == "" {
// 		return nil, errors.New("empty methodName in Encode")
// 	}
// 	if msg == nil {
// 		return nil, errors.New("nil message in Encode")
// 	}
// 	methodCall := &ProtobufMethodCall{
// 		Method: method,
// 		Msg:    msg,
// 	}
// 	return proto.Marshal(methodCall)
// }

func Encode(method string, msg proto.Message) ([]byte, error) {
	if method == "" {
		return nil, errors.New("empty method name in Proto Encode")
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	encodedMsg := &pb.EncodedMessage{
		Method: method,
		Data:   data,
	}

	encodedData, err := proto.Marshal(encodedMsg)
	if err != nil {
		return nil, err
	}

	return encodedData, nil
}

// func (p *ProtobufMessageCodec) Decode(method string, data []byte) (proto.Message, error) {
// 	if method == "" {
// 		return nil, errors.New("empty methodName in Decode")
// 	}
// 	if data == nil || len(data) == 0 {
// 		return nil, errors.New("empty data in Decode")
// 	}
// 	methodCall := &ProtobufMethodCall{}
// 	if err := proto.Unmarshal(data, methodCall); err != nil {
// 		return nil, err
// 	}
// 	return methodCall.Msg, nil
// }

func Decode(b []byte, msg proto.Message) error {
	return proto.Unmarshal(b, msg)
}
