package main

import (
	"bytes"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/pkg/term"
	"math"
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
// TODO use object for user for easier multiplayer
type user struct {
	position int
	character string
	mark string
}

	func (player *user) moveUser(d Direction){
	switch d {
	case Left:
		if player.position % 3 != 0 {
			player.position--
		}
	case Right:
		if player.position % 3 != 2 {
			player.position++
		}
	case Up:
		if math.Floor(float64(player.position / 3)) != 0 {
			player.position -= 3
		}
	case Down:
		if math.Floor(float64(player.position / 3)) != 2 {
			player.position += 3
		}
	default:
		println("Error has occurred in the move user function.")
	}
}

func (player user) placeMark(){
	originGrid[player.position] = player.mark
}

var originGrid [9]string
var displayGrid = make([]interface{}, len(originGrid))
var player1 = user{position: 0, character: "Y", mark: "Y"}

func gameLoop() {
	initGrid()
	for {
		draw()
		update()
		time.Sleep(1)
	}
}
func populateDisplayGrid(){
	for i, v := range originGrid {
		displayGrid[i] = v
	}
	displayGrid[player1.position] = player1.character
}
func update(){
	handleKeyEvents()
	populateDisplayGrid()
}

func draw(){
	tm.Clear() // Clear current screen
	tm.MoveCursor(1,1)
	tm.Printf(GRID_TEMPLATE, displayGrid...)
	tm.Flush()
}

func handleKeyEvents(){
	c := getch()
	switch {
	// TODO Add quit functionality
	case bytes.Equal(c, LEFT_KEY) || bytes.Equal(c, A_KEY): // left
		fmt.Println("LEFT pressed")
		player1.moveUser(Left)
	case bytes.Equal(c, RIGHT_KEY) || bytes.Equal(c, D_KEY): // right
		fmt.Println("RIGHT pressed")
		player1.moveUser(Right)
	case bytes.Equal(c, UP_KEY) || bytes.Equal(c, W_KEY): // up
		fmt.Println("UP pressed")
		player1.moveUser(Up)
	case bytes.Equal(c, DOWN_KEY) || bytes.Equal(c, S_KEY): // down
		fmt.Println("DOWN pressed")
		player1.moveUser(Down)
	case bytes.Equal(c, SPACE_KEY): // Place key
		fmt.Println("SPACE pressed")
		player1.placeMark()
	default:
		fmt.Println("Unknown pressed", c)
	}
}

func initGrid(){
	for i := 0; i < len(originGrid); i++ {
		originGrid[i] = "."
	}
	populateDisplayGrid()
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
