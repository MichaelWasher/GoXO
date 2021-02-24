package game

import "github.com/MichaelWasher/GoXO/grpc"

// User represents a single player within the game grid
type User struct {
	Position  int
	Character string
	Mark      string
	Name string
	InputChannel *chan grpc.Move
}
