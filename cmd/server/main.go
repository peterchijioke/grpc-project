package main

import (
	"os"
	"log"

	"grpc-project/internal/server"
)

func main() {
	port := ":50051"
	if p := os.Getenv("GRPC_PORT"); p != "" {
		port = ":" + p
	}

	if err := server.StartServer(port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
