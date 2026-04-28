package server

import (
	"context"
	"log"
	"net"

	"grpc-project/api/proto"
	"grpc-project/internal/auth"
	"grpc-project/internal/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

