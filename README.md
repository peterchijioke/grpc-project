# Go gRPC Project

A production-ready Go gRPC application with authentication (signup/login), unary and streaming RPC examples. Uses PostgreSQL with GORM and JWT authentication.

## Table of Contents

- [Project Structure](#project-structure)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running with Docker](#running-with-docker)
- [Running Locally](#running-locally)
- [API Reference](#api-reference)
- [Testing](#testing)
- [Adding New Services](#adding-new-services)

## Project Structure

```
.
├── api/
│   └── proto/                   # Protocol Buffer definitions
│       ├── greeter.proto
│       ├── auth.proto
│       ├── greeter.pb.go           # Generated Go code
│       ├── greeter_grpc.pb.go
│       ├── auth.pb.go
│       └── auth_grpc.pb.go
├── cmd/
│   ├── server/                  # Server entry point
│   │   └── main.go
│   └── client/                  # Client entry point
│       └── main.go
├── internal/
│   ├── server/                  # Server orchestration
│   │   └── server.go
│   ├── auth/                    # Authentication service
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── client/                  # Client helpers
│   │   └── client.go
│   ├── db/                      # Database connection
│   │   └── db.go
│   └── store/                   # Data models
│       └── user.go
├── configs/                     # Configuration files
├── scripts/                     # Helper scripts
│   └── init.sql                 # DB init SQL
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```
.
├── api/
│   └── proto/                   # Protocol Buffer definitions
│       ├── greeter.proto
│       ├── auth.proto
│       ├── greeter.pb.go           # Generated Go code
│       ├── greeter_grpc.pb.go      # Generated gRPC code
│       ├── auth.pb.go
│       └── auth_grpc.pb.go
├── cmd/
│   ├── server/                  # Server entry point
│   │   └── main.go
│   └── client/                  # Client entry point
│       └── main.go
├── internal/
│   ├── server/                  # Server implementation
│   │   └── server.go
│   ├── client/                  # Client implementation
│   │   └── client.go
│   ├── auth/                    # Auth service
│   │   └── auth.go
│   ├── db/                      # Database connection
│   │   └── db.go
│   └── store/                   # Data models
│       └── user.go
├── configs/                     # Configuration files
├── scripts/                     # Helper scripts
│   └── init.sql                 # Database initialization
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .env.example
├── .gitignore
└── README.md
```

## Features

- **Greeter Service**: Unary and server-streaming RPCs
- **Auth Service**: User signup and login with JWT tokens
- **PostgreSQL**: Persistent user storage with GORM ORM
- **Password Security**: Bcrypt password hashing
- **JWT Authentication**: Token-based auth (24h expiry)
- **Docker Support**: Full containerization with docker-compose
- **Reflection**: gRPC server reflection for grpcurl

## Prerequisites

- Go 1.25 or higher
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

## Configuration

Environment variables (use `.env` file for local dev):

| Variable      | Default               | Description              |
|---------------|----------------------|--------------------------|
| `GRPC_PORT`   | `50051`              | gRPC server port         |
| `DB_HOST`     | `localhost`          | PostgreSQL host          |
| `DB_PORT`     | `5432`               | PostgreSQL port          |
| `DB_USER`     | `postgres`           | PostgreSQL user          |
| `DB_PASSWORD` | `password`           | PostgreSQL password      |
| `DB_NAME`     | `grpc_auth`          | PostgreSQL database name |
| `DB_SSLMODE`  | `disable`            | SSL mode                 |
| `DB_DEBUG`    | `false`              | Enable GORM logging      |
| `JWT_SECRET`  | `super-secret-key-change-in-production` | JWT signing secret |

Copy `.env.example` to `.env` and adjust values.

## Generating gRPC Code

Regenerate Go code from `.proto` files:

```bash
make generate
# or manually:
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       api/proto/*.proto
```

## Running with Docker (Recommended)

### Using Docker Compose (includes PostgreSQL)

```bash
# Start all services (PostgreSQL + gRPC server)
make docker-compose-up

# Check logs
docker-compose logs -f grpc-server

# Stop services
make docker-compose-down
```

Services start on:
- gRPC server: `localhost:50051`
- PostgreSQL: `localhost:5432`

### Build and Run Server Only

```bash
# Build image
make docker-build

# Run container
make docker-run
```

## Running Locally

### 1. Start PostgreSQL

Using Docker:

```bash
docker run -d \
  --name grpc-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=grpc_auth \
  -p 5432:5432 \
  postgres:15-alpine
```

Or install locally via `brew install postgresql` and create database:

```bash
createdb grpc_auth
```

### 2. Run the Server

```bash
# With auto-reload (if needed)
go run ./cmd/server

# Or build and run
make build-server
./bin/server
```

### 3. Use the Client

#### Signup a new user

```bash
make run-client CMD=signup EMAIL=alice@example.com PASSWORD=secret123 NAME="Alice"
```

Or manually:

```bash
go run cmd/client/main.go \
  -cmd signup \
  -email alice@example.com \
  -password secret123 \
  -name "Alice"
```

#### Login

```bash
make run-client CMD=login EMAIL=alice@example.com PASSWORD=secret123
```

#### Greet (unary)

```bash
make run-client CMD=greet NAME=Bob
```

#### Greet (stream)

```bash
make run-client CMD=stream NAME=Charlie
```

## API Reference

### Auth Service

#### Signup

Creates a new user account and returns a JWT token.

**Request:**
```proto
message SignupRequest {
  string email = 1;
  string password = 2;
  string name = 3;
}
```

**Response:**
```proto
message SignupResponse {
  string user_id = 1;
  string email = 2;
  string name = 3;
  string token = 4;
  google.protobuf.Timestamp created_at = 5;
}
```

#### Login

Authenticates a user and returns a JWT token.

**Request:**
```proto
message LoginRequest {
  string email = 1;
  string password = 2;
}
```

**Response:**
```proto
message LoginResponse {
  string user_id = 1;
  string email = 2;
  string name = 3;
  string token = 4;
}
```

### Greeter Service

#### SayHello (Unary)

Simple unary RPC that returns a greeting.

#### SayHelloStream (Server Streaming)

Returns a stream of greeting messages.

## Testing with grpcurl

Install grpcurl:

```bash
brew install grpcurl
# or
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

List services:

```bash
grpcurl -plaintext localhost:50051 list
```

Call Signup:

```bash
grpcurl -plaintext -d '{"email":"test@example.com","password":"secret123","name":"Test User"}' \
  localhost:50051 auth.Auth/Signup
```

Call Login:

```bash
grpcurl -plaintext -d '{"email":"test@example.com","password":"secret123"}' \
  localhost:50051 auth.Auth/Login
```

Call SayHello:

```bash
grpcurl -plaintext -d '{"name":"World"}' \
  localhost:50051 greeter.Greeter/SayHello
```

## Running Tests

```bash
make test
# or
go test -v ./...
```

## Adding New Services

1. Create a new `.proto` file in `api/proto/`
2. Define service and messages
3. Regenerate code: `make generate`
4. Implement server logic in `internal/` (e.g., `internal/myservice/`)
5. Register service in `internal/server/server.go`
6. Add client methods in `internal/client/client.go`

## Security Notes

- **Never commit secrets**. Use `.env` file (already in `.gitignore`)
- Change `JWT_SECRET` in production (use strong random key)
- Use SSL/TLS in production (`grpc.WithTransportCredentials`)
- Enable `DB_DEBUG=false` in production
- Consider adding rate limiting and validation middleware

## Docker Commands

```bash
# View logs
docker-compose logs -f

# Access PostgreSQL shell
docker-compose exec postgres psql -U postgres -d grpc_auth

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v

# Rebuild server image
docker-compose build --no-cache grpc-server
```

## Troubleshooting

**Database connection errors:**
- Ensure PostgreSQL is running: `docker-compose ps`
- Check credentials in `.env` or docker-compose.yml
- DB might need time to start; docker-compose healthcheck handles this

**Port already in use:**
```bash
# Change port in .env or docker-compose.yml
export GRPC_PORT=50052
```

**Proto generation errors:**
```bash
# Ensure protoc plugins are installed
protoc --version
which protoc-gen-go
which protoc-gen-go-grpc
```

## License

MIT
