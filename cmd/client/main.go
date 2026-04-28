package main

import (
	"flag"
	"log"

	"grpc-project/internal/client"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "server address")
	command := flag.String("cmd", "greet", "command: signup, login, greet, stream")
	email := flag.String("email", "", "email address")
	password := flag.String("password", "", "password")
	name := flag.String("name", "", "user name")
	flag.Parse()

	client, err := client.Connect(*addr)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	switch *command {
	case "signup":
		if *email == "" || *password == "" || *name == "" {
			log.Fatal("email, password, and name are required for signup")
		}
		res, err := client.Signup(*email, *password, *name)
		if err != nil {
			log.Fatalf("signup failed: %v", err)
		}
		log.Printf("Signup successful! User ID: %s, Email: %s, Name: %s, Token: %s",
			res.UserId, res.Email, res.Name, res.Token)

	case "login":
		if *email == "" || *password == "" {
			log.Fatal("email and password are required for login")
		}
		res, err := client.Login(*email, *password)
		if err != nil {
			log.Fatalf("login failed: %v", err)
		}
		log.Printf("Login successful! User ID: %s, Email: %s, Name: %s, Token: %s",
			res.UserId, res.Email, res.Name, res.Token)

	case "greet":
		name := *name
		if name == "" {
			name = "World"
		}
		if err := client.SayHello(name); err != nil {
			log.Fatalf("greet failed: %v", err)
		}

	case "stream":
		name := *name
		if name == "" {
			name = "World"
		}
		if err := client.SayHelloStream(name); err != nil {
			log.Fatalf("stream failed: %v", err)
		}

	default:
		log.Fatal("unknown command. use 'signup', 'login', 'greet', or 'stream'")
	}
}
