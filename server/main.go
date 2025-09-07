package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "example.com/todo-grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type memStore struct {
	mu    sync.RWMutex
	next  int64
	items map[int64]*pb.Todo
}

func newStore() *memStore {
	return &memStore{next: 1, items: make(map[int64]*pb.Todo)}
}

type server struct {
	pb.UnimplementedTodoServiceServer
	s *memStore
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title required")
	}
	s.s.mu.Lock()
	defer s.s.mu.Unlock()
	id := s.s.next
	s.s.next++
	item := &pb.Todo{Id: id, Title: req.Title, Done: false}
	s.s.items[id] = item
	return &pb.CreateResponse{Item: item}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.s.mu.RLock()
	defer s.s.mu.RUnlock()
	it, ok := s.s.items[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return &pb.GetResponse{Item: it}, nil
}

func (s *server) List(req *pb.ListRequest, stream pb.TodoService_ListServer) error {
	s.s.mu.RLock()
	defer s.s.mu.RUnlock()
	for _, it := range s.s.items {
		time.Sleep(50 * time.Millisecond)
		if err := stream.Send(it); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, &server{s: newStore()})

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
