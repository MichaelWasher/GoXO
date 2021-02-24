package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func SetupClient(portNum int, inputChannel *chan Move){

	log.Println("Configuring the Client")

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%d", portNum), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewMoveServiceClient(conn)
	
	//TODO check int
	for {
		RequestBody := <- *inputChannel

		_, err = c.PerformMove(context.Background(), &MoveMessage{MoveRequest: RequestBody})

		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}

		// Break on quit
		if RequestBody == Move_Quit{
			break
		}
	}
}