package io

// TODO Add Thread Safety
import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/pkg/term"
)

// ---- Constants
type KeyMap struct {
	Down      []byte
	Up        []byte
	Right     []byte
	Left      []byte
	PlaceMark []byte
	Quit      []byte
}

var ArrowKeyMap = KeyMap{
	Down:      []byte{27, 91, 66}, // Down Arrow
	Up:        []byte{27, 91, 65}, // Up Arrow
	Right:     []byte{27, 91, 67}, // Right Arrow
	Left:      []byte{27, 91, 68}, // Left Arrow
	PlaceMark: []byte{32},         //Space
	Quit:      []byte{113},        //Q
}
var AlphaKeyMap = KeyMap{
	Down:      []byte{115}, // S
	Up:        []byte{119}, // W
	Right:     []byte{100}, // D
	Left:      []byte{97},  // A
	PlaceMark: []byte{112}, //P
	Quit:      []byte{3},   //CTRL + C
}

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
			// TODO Configure Multiple User Keymaps
			playerInput <- getInput(getch(&t.Term), ArrowKeyMap)

		case <-ctx.Done():
			return
		}
	}
}

func getInput(c []byte, km KeyMap) InputEvent {

	switch {
	case bytes.Equal(c, km.Left):
		return NewInputEvent(Move_Left)

	case bytes.Equal(c, km.Right):
		return NewInputEvent(Move_Right)

	case bytes.Equal(c, km.Up):
		return NewInputEvent(Move_Up)

	case bytes.Equal(c, km.Down):
		return NewInputEvent(Move_Down)

	case bytes.Equal(c, km.PlaceMark):
		return NewInputEvent(Move_PlaceMark)

	case bytes.Equal(c, km.Quit):
		tmpEvent := NewInputEvent(Move_Quit)
		tmpEvent.Terminate = true
		return tmpEvent

	default:
		// TODO Mute this
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
