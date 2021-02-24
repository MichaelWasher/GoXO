package grpc

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// Server represents the gRPC server
type Server struct {
}

var outputChannel *chan Move

// PerformMove generates response to a Ping request
func (s *Server) PerformMove(ctx context.Context, in *MoveMessage) (*Empty, error) {
	log.Printf("Receive message %s", in.MoveRequest)
	// Write to the Channel //
	*outputChannel <- in.MoveRequest
	return &Empty{}, nil
}

// TODO Rewrite this to be clean
func SetupServer(portNum int, channel *chan Move){
	outputChannel = channel

	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portNum))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	s := Server{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	reflection.Register(grpcServer)
	RegisterMoveServiceServer(grpcServer, &s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}