package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-project/api/proto"
)

type AuthClient interface {
	Signup(email, password, name string) (*proto.SignupResponse, error)
	Login(email, password string) (*proto.LoginResponse, error)
}

type GreeterClient interface {
	SayHello(name string) error
	SayHelloStream(name string) error
}

type Client struct {
	conn    *grpc.ClientConn
	greeter proto.GreeterClient
	auth    proto.AuthClient
}

func Connect(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		greeter: proto.NewGreeterClient(conn),
		auth:    proto.NewAuthClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Greeter methods
func (c *Client) SayHello(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.greeter.SayHello(ctx, &proto.HelloRequest{Name: name})
	if err != nil {
		return err
	}
	log.Printf("Response: %s", res.GetMessage())
	return nil
}

func (c *Client) SayHelloStream(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.greeter.SayHelloStream(ctx, &proto.HelloRequest{Name: name})
	if err != nil {
		return err
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("Stream message: %s", res.GetMessage())
	}
	return nil
}

// Auth methods
func (c *Client) Signup(email, password, name string) (*proto.SignupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return c.auth.Signup(ctx, &proto.SignupRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})
}

func (c *Client) Login(email, password string) (*proto.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return c.auth.Login(ctx, &proto.LoginRequest{
		Email:    email,
		Password: password,
	})
}
