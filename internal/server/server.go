package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"grpc-project/api/proto"
	"grpc-project/internal/auth"
	"grpc-project/internal/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GreeterServer struct {
	proto.UnimplementedGreeterServer
}

type AuthServer struct {
	proto.UnimplementedAuthServer
	authService *auth.AuthService
}

func NewAuthServer() *AuthServer {
	return &AuthServer{
		authService: auth.NewAuthService(),
	}
}

func (s *GreeterServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	if err := validateHelloRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	name := strings.TrimSpace(req.GetName())
	log.Printf("Received request from: %s", name)
	return &proto.HelloReply{
		Message: "Hello, " + name + "!",
	}, nil
}

func (s *GreeterServer) SayHelloStream(req *proto.HelloRequest, stream proto.Greeter_SayHelloStreamServer) error {
	if err := validateHelloRequest(req); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	name := strings.TrimSpace(req.GetName())
	log.Printf("Stream request from: %s", name)
	messages := []string{
		"Hello, " + name + "!",
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

func validateHelloRequest(req *proto.HelloRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) > 100 {
		return fmt.Errorf("name must be less than 100 characters")
	}

	return nil
}

func (as *AuthServer) Signup(ctx context.Context, req *proto.SignupRequest) (*proto.SignupResponse, error) {
	return as.authService.Signup(ctx, req)
}

func (as *AuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return as.authService.Login(ctx, req)
}

func StartServer(port string) error {
	// Connect to database
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &GreeterServer{})
	proto.RegisterAuthServer(s, NewAuthServer())

	// Enable reflection for grpcurl
	reflection.Register(s)

	log.Printf("Server listening on %s", port)
	return s.Serve(lis)
}
