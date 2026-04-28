package server

import (
	"context"
	"log"
	"net"

	"grpc-project/api/proto"
	"google.golang.org/grpc"
)

type GreeterServer struct {
	proto.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	log.Printf("Received request from: %s", req.GetName())
	return &proto.HelloReply{
		Message: "Hello, " + req.GetName() + "!",
	}, nil
}

func (s *GreeterServer) SayHelloStream(req *proto.HelloRequest, stream proto.Greeter_SayHelloStreamServer) error {
	log.Printf("Stream request from: %s", req.GetName())
	messages := []string{
		"Hello, " + req.GetName() + "!",
		"Nice to meet you!",
		"Welcome to gRPC!",
		"Goodbye!",
	}

	for _, msg := range messages {
		if err := stream.Send(&proto.HelloReply{Message: msg}); err != nil {
			return err
		}
	}
	return nil
}

func StartServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &GreeterServer{})

	log.Printf("Server listening on %s", port)
	return s.Serve(lis)
}
