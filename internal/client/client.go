package client

import (
	"context"
	"log"
	"time"

	"grpc-project/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Connect(addr string) (proto.GreeterClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return proto.NewGreeterClient(conn), nil
}

func SayHello(client proto.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SayHello(ctx, &proto.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %s", res.GetMessage())
}

func SayHelloStream(client proto.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.SayHelloStream(ctx, &proto.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not start stream: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("Stream message: %s", res.GetMessage())
	}
}
