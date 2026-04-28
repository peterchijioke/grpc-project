package main

import (
	"flag"
	"log"

	"grpc-project/internal/client"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "server address")
	mode := flag.String("mode", "unary", "mode: unary or stream")
	name := flag.String("name", "World", "name to greet")
	flag.Parse()

	greeterClient, err := client.Connect(*addr)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	switch *mode {
	case "unary":
		client.SayHello(greeterClient, *name)
	case "stream":
		client.SayHelloStream(greeterClient, *name)
	default:
		log.Fatal("unknown mode. use 'unary' or 'stream'")
	}
}
