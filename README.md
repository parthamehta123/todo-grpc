# gRPC Todo Example (Go)

A simple gRPC service for managing todos, built in Go with Protocol Buffers.  
This project demonstrates defining a proto service, generating Go code, and implementing both server and client.

---

## ✨ Features

- Protobuf definition (`api/todo.proto`)
- gRPC service with methods:
  - `Create` → Create a todo
  - `Get` → Fetch a todo by ID
  - `List` → Stream all todos (server streaming)
- In-memory store with RWMutex
- Separate server and client implementations
- Makefile for easy build/run/gen

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- `protoc` installed
- `protoc-gen-go` and `protoc-gen-go-grpc` installed
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

### Clone the repo
```bash
git clone https://github.com/YOUR_USERNAME/todo-grpc.git
cd todo-grpc
```

### Generate protobuf code
```bash
make gen
```

This generates `api/todo.pb.go` and `api/todo_grpc.pb.go`.

---

## 📚 Example Usage

### Run server
```bash
make run-server
```

### Run client (in another terminal)
```bash
make run-client
```

### Sample output
```
created: id:1  title:"learn gRPC"
fetched: id:1  title:"learn gRPC"
stream item: id:1  title:"learn gRPC"
stream end: EOF
```

---

## 🧪 Development

Run tests:
```bash
go test ./...
```

Build manually:
```bash
go build -o bin/server ./server
go build -o bin/client ./client
```

---

## 📖 Notes

- Data is **not persisted** (in-memory only).
- A good starting point for learning gRPC in Go.
- Can be extended with a real database, authentication, or bidirectional streaming.
