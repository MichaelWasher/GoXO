package main

import (
	"bytes"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/pkg/term"
	"time"
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
const GRID_TEMPLATE = `
-------------
| %s | %s | %s |
-------------
| %s | %s | %s |
-------------
| %s | %s | %s |
-------------
`

// ---- Game Variables
// TODO use object for user for easier multiplayer
var userPosition = 0
var grid [9]string
var userCharacter string = "Y"

func gameLoop() {
	initGrid()
	for {
		draw()
		update()
		time.Sleep(1)
	}
}

func update(){
	handleKeyEvents()
}

func draw(){
	tm.Clear() // Clear current screen
	tm.MoveCursor(1,1)
	s := prepareGridForPrint(grid)
	tm.Printf(GRID_TEMPLATE, s...)
	tm.Flush()
}

func handleKeyEvents(){
	c := getch()
	switch {
	// TODO Add quit functionality
	case bytes.Equal(c, LEFT_KEY) || bytes.Equal(c, A_KEY): // left
		fmt.Println("LEFT pressed")
	case bytes.Equal(c, RIGHT_KEY) || bytes.Equal(c, D_KEY): // right
		fmt.Println("RIGHT pressed")
		userPosition++
	case bytes.Equal(c, UP_KEY) || bytes.Equal(c, W_KEY): // up
		fmt.Println("UP pressed")
	case bytes.Equal(c, DOWN_KEY) || bytes.Equal(c, S_KEY): // down
		fmt.Println("DOWN pressed")
	case bytes.Equal(c, SPACE_KEY): // Place key
		fmt.Println("SPACE pressed")
	default:
		fmt.Println("Unknown pressed", c)
	}
}

func initGrid(){
	for i := 0; i < len(grid); i++ {
		grid[i] = "."
	}
}
func prepareGridForPrint(grid [9] string) []interface{} {
	s := make([]interface{}, len(grid))
	for i, v := range grid {
		s[i] = v
	}
	s[userPosition] = userCharacter
	return s
}

// ---- Utility Functions
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
