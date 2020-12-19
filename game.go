package main

import (
	"bytes"
	"fmt"
	"time"

	tm "github.com/buger/goterm"
	"github.com/pkg/term"
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

type Direction int

const (
	Left Direction = 1 << iota
	Right
	Up
	Down
)

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
var originGrid [9]string
var displayGrid = make([]interface{}, len(originGrid))

// Users
var player1 = User{Position: 0, Character: "Y", Mark: "X"}
var player2 = User{Position: 8, Character: "Y", Mark: "0"}
var running bool

// Core Game Loop

func gameLoop() {
	initGrid()
	running = true
	for running {
		draw()
		update()
		time.Sleep(1)
	}
}
func populateDisplayGrid() {
	for i, v := range originGrid {
		displayGrid[i] = v
	}
	displayGrid[player1.Position] = player1.Character
}
func update() {
	handleKeyEvents()
	populateDisplayGrid()
}

func draw() {
	tm.Clear() // Clear current screen
	tm.MoveCursor(1, 1)
	tm.Printf(GRID_TEMPLATE, displayGrid...)
	tm.Flush()
}

func handleKeyEvents() {
	c := getch()
	for {
		switch {
		// TODO Add quit functionality
		case bytes.Equal(c, LEFT_KEY) || bytes.Equal(c, A_KEY): // left
			fmt.Println("LEFT pressed")
			player1.MoveUser(Left)
		case bytes.Equal(c, RIGHT_KEY) || bytes.Equal(c, D_KEY): // right
			fmt.Println("RIGHT pressed")
			player1.MoveUser(Right)
		case bytes.Equal(c, UP_KEY) || bytes.Equal(c, W_KEY): // up
			fmt.Println("UP pressed")
			player1.MoveUser(Up)
		case bytes.Equal(c, DOWN_KEY) || bytes.Equal(c, S_KEY): // down
			fmt.Println("DOWN pressed")
			player1.MoveUser(Down)
		case bytes.Equal(c, SPACE_KEY): // Place key
			fmt.Println("SPACE pressed")
			player1.PlaceMark()
			break
		case bytes.Equal(c, Q_KEY):
			running = false
			break
		default:
			fmt.Println("Unknown pressed", c)
		}
		if originGrid[player1.Position] == "." {
			break
		}
	}
}

func initGrid() {
	for i := 0; i < len(originGrid); i++ {
		originGrid[i] = "."
	}
	populateDisplayGrid()
}

// TODO -- Implement the game logic. Wining, Losing and Points
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
