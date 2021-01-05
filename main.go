package main

import (
	"fmt"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// TODO Add Multiplayer Support
// TODO Add Socket Support for Multiple Player Input
// TODO Flags for the Socket Connection
func main() {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	s := Server{}
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	//reflection.Register(&s)
	RegisterRequestMoveServer(grpcServer, &s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	//gameLoop()
}
