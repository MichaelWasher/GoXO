package io

import (
	"bytes"
	"context"
	"fmt"
	"log"

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
var CTRL_C_KEYS = []byte{3}

// Terminal Definition - Matches the IO Handler Interface
type Terminal struct {
	term.Term
}

func NewTerminal() Terminal {
	// Create Terminal
	terminal := Terminal{}
	t, err := term.Open("/dev/tty")
	if err != nil {
		log.Fatalf("Unable to open a terminal. %v", err)
	}
	terminal.Term = *t

	// Configure Terminal into Raw Mode
	err = t.SetRaw()
	if err != nil {
		log.Fatalf("Unable to set Terminal into Raw Mode. %v", err)
		log.Fatal("Attempting to continue but the game may not display correctly.")
	}

	return terminal
}

func (t Terminal) Close() {
	// Set terminal back to sane mode
	defer t.Term.Close()
	t.SetCbreak()
}

func (t Terminal) Print(outString string) {
	t.Write([]byte(outString))
	t.Flush()
}

func (t Terminal) RegisterDrawEvents(ctx context.Context, drawChannel <-chan DrawEvent) {
	// While Not Quit
	for true {
		select {
		case <-ctx.Done():
			return
		case event := <-drawChannel:
			// Clear current screen
			// TODO this does not have cross-platform support
			// TODO This does not respect the current TTY access. Should print to the TTY but this has buffering issues.
			fmt.Print("\033[H\033[2J")
			fmt.Print(event.DrawString)
		}

	}
}

func (t Terminal) RegisterInputEvents(ctx context.Context, playerInput chan InputEvent) {
	for {
		select {
		case <-playerInput:
			playerInput <- getInput(getch(&t.Term))

		case <-ctx.Done():
			return
		}
	}
}

func getInput(c []byte) InputEvent {

	switch {
	case bytes.Equal(c, LEFT_KEY) || bytes.Equal(c, A_KEY): // left
		log.Print("LEFT pressed")
		return NewInputEvent(Move_Left)

	case bytes.Equal(c, RIGHT_KEY) || bytes.Equal(c, D_KEY): // right
		log.Print("RIGHT pressed")
		return NewInputEvent(Move_Right)

	case bytes.Equal(c, UP_KEY) || bytes.Equal(c, W_KEY): // up
		log.Print("UP pressed")
		return NewInputEvent(Move_Up)

	case bytes.Equal(c, DOWN_KEY) || bytes.Equal(c, S_KEY): // down
		log.Print("DOWN pressed")
		return NewInputEvent(Move_Down)

	case bytes.Equal(c, SPACE_KEY): // Place key
		log.Print("SPACE pressed")
		return NewInputEvent(Move_PlaceMark)

	case bytes.Equal(c, Q_KEY) || bytes.Equal(c, CTRL_C_KEYS):
		log.Print("Q pressed")
		tmpEvent := NewInputEvent(Move_Quit)
		tmpEvent.Terminate = true
		return tmpEvent

	default:
		log.Println("Unknown pressed", c)
	}
	return NewInputEvent(Move_Noop)
}

func getch(terminal *term.Term) []byte {

	bytes := make([]byte, 3)
	numRead, err := terminal.Read(bytes)
	if err != nil {
		log.Printf("Unable to read from the Terminal. %v", err)
		return nil
	}

	return bytes[0:numRead]

}
