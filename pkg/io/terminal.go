package io

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

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

func NewTerminal() (*Terminal, error) {
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
	err = t.SetReadTimeout(500 * time.Millisecond)
	if err != nil {
		log.Fatalf("Unable to configure the Terminal read timeout.", err)
		return nil, errors.New("unable to set terminal read timeout")
	}

	return &terminal, nil
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

func (t *Terminal) RegisterDrawEvents(ctx context.Context, drawChannel <-chan DrawEvent) {
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

func (t *Terminal) RegisterInputEvents(ctx context.Context, playerInput chan InputEvent) {
	charInputChannel := t.getCharacterInputChannel(ctx)
	for {
		select {
		case char := <-charInputChannel:
			playerInput <- inputEventFromBytes(char)
		case <-ctx.Done():
			return
		}
	}
}

func inputEventFromBytes(c []byte) InputEvent {

	switch {
	case bytes.Equal(c, LEFT_KEY) || bytes.Equal(c, A_KEY): // left
		return NewInputEvent(Move_Left)

	case bytes.Equal(c, RIGHT_KEY) || bytes.Equal(c, D_KEY): // right
		return NewInputEvent(Move_Right)

	case bytes.Equal(c, UP_KEY) || bytes.Equal(c, W_KEY): // up
		return NewInputEvent(Move_Up)

	case bytes.Equal(c, DOWN_KEY) || bytes.Equal(c, S_KEY): // down
		return NewInputEvent(Move_Down)

	case bytes.Equal(c, SPACE_KEY): // Place key
		return NewInputEvent(Move_PlaceMark)

	case bytes.Equal(c, Q_KEY) || bytes.Equal(c, CTRL_C_KEYS):
		tmpEvent := NewInputEvent(Move_Quit)
		tmpEvent.Terminate = true
		return tmpEvent

	default:
		log.Println("Unknown pressed", c)
	}
	return NewInputEvent(Move_Noop)
}

func (t *Terminal) getCharacterInputChannel(ctx context.Context) chan []byte {
	characterInputChannel := make(chan []byte)
	go func() {
		for {

			select {
			case <-ctx.Done():
				return

			default:
				bytes := make([]byte, 3)
				numRead, err := t.Read(bytes)

				// Did not read anything from terminal in time (No need to log as this is expected)
				if err == io.EOF {
					continue
				}

				if err != nil {
					log.Printf("Unable to read from the Terminal. %v", err)
					continue
				}

				log.Print("Collected bytes below:")
				log.Print(bytes)

				characterInputChannel <- bytes[0:numRead]
			}
		}

	}()
	return characterInputChannel
}
