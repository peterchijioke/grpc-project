package server

import (
	"context"
	"testing"

	"grpc-project/api/proto"
)

func TestSayHello(t *testing.T) {
	s := &GreeterServer{}
	req := &proto.HelloRequest{Name: "Test"}

	res, err := s.SayHello(context.Background(), req)
	if err != nil {
		t.Fatalf("SayHello() error = %v", err)
	}

	expected := "Hello, Test!"
	if res.GetMessage() != expected {
		t.Errorf("SayHello() = %v, want %v", res.GetMessage(), expected)
	}
}

func TestSayHello_EmptyName(t *testing.T) {
	s := &GreeterServer{}
	req := &proto.HelloRequest{Name: ""}

	res, err := s.SayHello(context.Background(), req)
	if err != nil {
		t.Fatalf("SayHello() error = %v", err)
	}

	expected := "Hello, !"
	if res.GetMessage() != expected {
		t.Errorf("SayHello() = %v, want %v", res.GetMessage(), expected)
	}
}
