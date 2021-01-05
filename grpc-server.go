package main

import (
	"golang.org/x/net/context"
	"log"
)

// Server represents the gRPC server
type Server struct {
}

// PerformMove generates response to a Ping request
func (s *Server) PerformMove(ctx context.Context, in *MoveMessage) (*Empty, error) {
	log.Printf("Receive message %s", in.MoveRequest)
	// Write to the Channel //
	outstandingMoves <- in.MoveRequest
	return &Empty{}, nil
}