package game

import "github.com/MichaelWasher/GoXO/pkg/io"

// User represents a single player within the game grid
type User struct {
	Position           int
	Character          string
	Mark               string
	Name               string
	InputChannel       chan io.InputEvent
	PlayerEventHandler io.PlayerInputHandler
}
