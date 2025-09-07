package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example.com/todo-grpc/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewTodoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create
	cr, err := c.Create(ctx, &pb.CreateRequest{Title: "learn gRPC"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("created:", cr.Item)

	// Get
	gr, err := c.Get(ctx, &pb.GetRequest{Id: cr.Item.Id})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("fetched:", gr.Item)

	// List (server streaming)
	stream, err := c.List(ctx, &pb.ListRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		item, err := stream.Recv()
		if err != nil {
			fmt.Println("stream end:", err)
			return
		}
		fmt.Println("stream item:", item)
	}
}
