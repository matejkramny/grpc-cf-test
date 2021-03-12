package main

import (
	"fmt"
	"log"
	"net"

	testgrpc "grpc_test/grpc"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":2015")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	testgrpc.RegisterTestServiceServer(srv, newServer())

	fmt.Println("Listening to :2015")
	srv.Serve(lis)
}
