package main

import (
	"bytes"
	"fmt"
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
			outstandingMoves <- Move_Left
		case bytes.Equal(c, RIGHT_KEY) || bytes.Equal(c, D_KEY): // right
			log.Print("RIGHT pressed")
			outstandingMoves <- Move_Right
		case bytes.Equal(c, UP_KEY) || bytes.Equal(c, W_KEY): // up
			log.Print("UP pressed")
			outstandingMoves <- Move_Up
		case bytes.Equal(c, DOWN_KEY) || bytes.Equal(c, S_KEY): // down
			log.Print("DOWN pressed")
			outstandingMoves <- Move_Down
		case bytes.Equal(c, SPACE_KEY): // Place key
			log.Print("SPACE pressed")
			outstandingMoves <- Move_PlaceMark
		case bytes.Equal(c, Q_KEY) || bytes.Equal(c, CTRL_C_KEYS):
			outstandingMoves <- Move_Quit
			break
		default:
			fmt.Println("Unknown pressed", c)
			continue
		}
	}
}

func getch() []byte {
	bytes := make([]byte, 3)
	numRead, err := terminal.Read(bytes)

	if err != nil {
		return nil
	}
	return bytes[0:numRead]
}
