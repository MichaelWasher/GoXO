package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/term"
	"log"
)

// ---- Constants
var DOWN_KEY = []byte{27, 91, 66}
var UP_KEY = []byte{27, 91, 65}
var RIGHT_KEY = []byte{27, 91, 67}
var LEFT_KEY = []byte{27, 91, 68}
var W_KEY = []byte{119}
var A_KEY = []byte{97}
var S_KEY = []byte{115}
var D_KEY = []byte{100}
var SPACE_KEY = []byte{32}
var Q_KEY = []byte{113}
var CTRL_C_KEYS = []byte{3}

// ---- Handle the Key events and hand off to the game routine
func handleKeyEvents() {
	for(running) {
		c := getch()

		switch {
		// TODO Add quit functionality
		case bytes.Equal(c, LEFT_KEY) || bytes.Equal(c, A_KEY): // left
			log.Print("LEFT pressed")
			outstandingMoves <- MoveLeft
		case bytes.Equal(c, RIGHT_KEY) || bytes.Equal(c, D_KEY): // right
			log.Print("RIGHT pressed")
			outstandingMoves <- MoveRight
		case bytes.Equal(c, UP_KEY) || bytes.Equal(c, W_KEY): // up
			log.Print("UP pressed")
			outstandingMoves <- MoveUp
		case bytes.Equal(c, DOWN_KEY) || bytes.Equal(c, S_KEY): // down
			log.Print("DOWN pressed")
			outstandingMoves <- MoveDown
		case bytes.Equal(c, SPACE_KEY): // Place key
			log.Print("SPACE pressed")
			outstandingMoves <- PlacePiece
		case bytes.Equal(c, Q_KEY) || bytes.Equal(c, CTRL_C_KEYS):
			outstandingMoves <- Quit
			break
		default:
			fmt.Println("Unknown pressed", c)
			continue
		}
		<- outstandingMoves
	}
}

func getch() []byte {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)
	numRead, err := t.Read(bytes)
	t.Restore()
	t.Close()
	if err != nil {
		return nil
	}
	return bytes[0:numRead]
}
