# Go gRPC Project

A standard Go gRPC application with unary and streaming RPC examples.

## Project Structure

```
.
├── api/
│   └── proto/               # Protocol Buffer definitions
│       ├── greeter.proto
│       ├── greeter.pb.go           # Generated Go code
│       └── greeter_grpc.pb.go      # Generated gRPC code
├── cmd/
│   ├── server/              # Server entry point
│   │   └── main.go
│   └── client/              # Client entry point
│       └── main.go
├── internal/
│   ├── server/              # Server implementation
│   │   └── server.go
│   └── client/              # Client implementation
│       └── client.go
├── configs/                 # Configuration files (optional)
├── scripts/                 # Helper scripts (optional)
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- Protocol Buffers compiler (`protoc`)
- Protocol Buffers Go plugins

### Install Dependencies

```bash
# Install protoc (macOS)
brew install protobuf

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install project dependencies
go mod download
```

## Generating gRPC Code

Regenerate Go code from `.proto` files:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       api/proto/*.proto
```

## Running the Server

```bash
# Run with default port 50051
go run cmd/server/main.go

# Or build and run
go build -o bin/server cmd/server/main.go
./bin/server
```

Configure port via environment variable:

```bash
export GRPC_PORT=8080
go run cmd/server/main.go
```

## Running the Client

### Unary RPC

```bash
go run cmd/client/main.go -addr localhost:50051 -mode unary -name "Alice"
```

### Server Streaming RPC

```bash
go run cmd/client/main.go -addr localhost:50051 -mode stream -name "Bob"
```

## Testing with curl (gRPCurl)

Install [grpcurl](https://github.com/fullstorydev/grpcurl) for testing:

```bash
# List services
grpcurl -plaintext localhost:50051 list

# Call SayHello method
grpcurl -plaintext -d '{"name":"Test"}' localhost:50051 greeter.Greeter/SayHello
```

## Docker Support

### Build and run with Docker Compose

```yaml
# docker-compose.yml (optional)
version: '3.8'
services:
  grpc-server:
    build: .
    ports:
      - "50051:50051"
```

Build and run:

```bash
docker build -t grpc-server .
docker run -p 50051:50051 grpc-server
```

## Adding New Services

1. Create a new `.proto` file in `api/proto/`
2. Define your service and messages
3. Regenerate code with `protoc`
4. Implement the server in `internal/server/`
5. Update the server main if needed

## Code Style

This project follows standard Go conventions:
- Clear package names
- Error handling with explicit checks
- Context-based cancellation
- Logging with the standard library

## Security Notes

- Never commit secrets or environment files
- Use `.env` files for local development (already in `.gitignore`)
- For production, use proper secret management (Vault, Kubernetes secrets, etc.)

## License

MIT
