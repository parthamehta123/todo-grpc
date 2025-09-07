PROTO=api/todo.proto

gen:
	protoc -I api --go_out=. --go-grpc_out=. $(PROTO)

run-server:
	go run ./server

run-client:
	go run ./client
