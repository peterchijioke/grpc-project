package server

import (
	"context"
	"strings"
	"testing"

	"grpc-project/api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	_, err := s.SayHello(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("unexpected error type: %v", err)
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected code InvalidArgument, got %v", st.Code())
	}
	if st.Message() != "name is required" {
		t.Errorf("expected message 'name is required', got %q", st.Message())
	}
}

func TestSayHello_LongName(t *testing.T) {
	s := &GreeterServer{}
	req := &proto.HelloRequest{Name: strings.Repeat("a", 101)}

	_, err := s.SayHello(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for long name, got nil")
	}

	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected code InvalidArgument, got %v", st.Code())
	}
}

func TestSayHello_WhitespaceName(t *testing.T) {
	s := &GreeterServer{}
	req := &proto.HelloRequest{Name: "   "}

	_, err := s.SayHello(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for whitespace-only name, got nil")
	}

	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected code InvalidArgument, got %v", st.Code())
	}
}
