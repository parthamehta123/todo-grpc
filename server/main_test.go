package main

import (
	"context"
	"net"
	"testing"
	"time"

	pb "example.com/todo-grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// Use an in-memory connection instead of a real TCP port
var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// spin up a test server with bufconn
func startTestServer(t *testing.T) pb.TodoServiceClient {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &server{s: newStore()})

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return pb.NewTodoServiceClient(conn)
}

func TestCreateAndGet(t *testing.T) {
	client := startTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create
	created, err := client.Create(ctx, &pb.CreateRequest{Title: "write tests"})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.Item.GetId() == 0 {
		t.Errorf("expected nonzero id")
	}
	if created.Item.GetTitle() != "write tests" {
		t.Errorf("expected title 'write tests', got %q", created.Item.GetTitle())
	}

	// Get
	got, err := client.Get(ctx, &pb.GetRequest{Id: created.Item.Id})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Item.GetId() != created.Item.Id {
		t.Errorf("expected id %d, got %d", created.Item.Id, got.Item.GetId())
	}
}

func TestList(t *testing.T) {
	client := startTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create two todos
	_, _ = client.Create(ctx, &pb.CreateRequest{Title: "task1"})
	_, _ = client.Create(ctx, &pb.CreateRequest{Title: "task2"})

	stream, err := client.List(ctx, &pb.ListRequest{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	var count int
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
		count++
	}
	if count < 2 {
		t.Errorf("expected at least 2 items, got %d", count)
	}
}
