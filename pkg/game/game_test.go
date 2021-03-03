package game_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/MichaelWasher/GoXO/pkg/game"
	"github.com/MichaelWasher/GoXO/pkg/io"
)

type MockGameIO struct {
	InputChannel chan io.InputEvent
}
type MockDrawIO struct {
	DrawChannel <-chan io.DrawEvent
}

func (mio *MockDrawIO) RegisterDrawEvents(ctx context.Context, drawChannel <-chan io.DrawEvent) {
	fmt.Print("Draw Event Called")
	mio.DrawChannel = drawChannel
	select {
	case <-ctx.Done():
		return
	}
}
func (mio *MockGameIO) RegisterInputEvents(ctx context.Context, inputChan chan io.InputEvent) {
	fmt.Print("Input Event Called")
	mio.InputChannel = inputChan
	select {
	case <-ctx.Done():
		return
	}
}
func (mio *MockGameIO) Write(ioe io.InputEvent) error {
	if mio.InputChannel == nil {
		return errors.New("Input Channel has not been configured correctly")
	}
	mio.InputChannel <- ioe
	return nil
}

func (mio MockDrawIO) Read() (io.DrawEvent, error) {
	if mio.DrawChannel == nil {
		return io.DrawEvent{}, errors.New("Unable to get draw event from IO tool. IO tool is nil")
	}
	return <-mio.DrawChannel, nil
}

func (mio MockGameIO) CloseGame() error {
	if mio.InputChannel == nil {
		return errors.New("Unable to close the game as there is no open Input Channel.")
	}
	mio.InputChannel <- io.NewInputEvent(io.Move_Quit)
	return nil

}
func SetupTest(t *testing.T) (*MockGameIO, *MockGameIO, *MockDrawIO, *game.Game) {
	// Configure the IO
	p1Mio := &MockGameIO{}
	p2Mio := &MockGameIO{}
	drawMio := &MockDrawIO{}

	// Create the Game
	gameObject := game.NewGame(p1Mio, p2Mio, drawMio)
	t.Log("Game Created")

	// Setup Game Loop
	go gameObject.GameLoop()
	// Wait for the Game to hook into the p1/p2/draw channels
	for {
		time.Sleep(10 * time.Millisecond)
		if drawMio.DrawChannel != nil {
			break
		}
	}

	return p1Mio, p2Mio, drawMio, gameObject
}
func TeardownTest(t *testing.T) {

}

// TODO Test Perform Win
// TODO Create Generic Test that can be used for iterative moves

func TestIO(t *testing.T) {
	// Test Wrapping Move Left
	testCases := []struct {
		name           string
		p1Moves        []io.Move
		p2Moves        []io.Move
		expectedOutput string
	}{
		{"Player 2 Move Up", []io.Move{io.Move_Left, io.Move_Down, io.Move_Right, io.Move_Up, io.Move_PlaceMark}, []io.Move{io.Move_Up}, "-------------\r\n| X | . | . |\r\n-------------\r\n| . | . | 2 |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Test
			p1Mio, p2Mio, drawMio, gameObject := SetupTest(t)
			defer gameObject.CloseGame()
			defer TeardownTest(t)

			// Discard the first read
			drawMio.Read()

			// Wait for users turn and perform movement
			var drawEvent io.DrawEvent
			var err error
			for _, move := range tc.p1Moves {
				t.Log("Input from P1 Mio")
				p1Mio.Write(io.InputEvent{Move: move, Terminate: false})

				// Read the Draw Output
				time.Sleep(5 * time.Millisecond)
				drawEvent, err = drawMio.Read()
				if err != nil {
					t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
				}
			}
			for _, move := range tc.p2Moves {
				t.Log("Input from P2 Mio")

				p2Mio.Write(io.InputEvent{Move: move, Terminate: false})
				// Read the Draw Output
				time.Sleep(5 * time.Millisecond)
				drawEvent, err = drawMio.Read()
				if err != nil {
					t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
				}
			}

			// Compare against the template
			if !strings.HasPrefix(drawEvent.DrawString, tc.expectedOutput) {
				t.Fatal("Moving Player 1 failed.")
			}
		})
	}
}

func TestMutiplayerTurns(t *testing.T) {
	// Test Wrapping Move Left
	testCases := []struct {
		name           string
		p1Moves        []io.Move
		p2Moves        []io.Move
		expectedOutput string
	}{
		{"Player 2 Move Up", []io.Move{io.Move_Left, io.Move_Down, io.Move_Right, io.Move_Up, io.Move_PlaceMark}, []io.Move{io.Move_Up}, "-------------\r\n| X | . | . |\r\n-------------\r\n| . | . | 2 |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Test
			p1Mio, p2Mio, drawMio, gameObject := SetupTest(t)
			defer gameObject.CloseGame()
			defer TeardownTest(t)

			// Discard the first read
			drawMio.Read()

			// Wait for users turn and perform movement
			var drawEvent io.DrawEvent
			var err error
			for _, move := range tc.p1Moves {
				t.Log("Input from P1 Mio")
				p1Mio.Write(io.InputEvent{Move: move, Terminate: false})

				// Read the Draw Output
				time.Sleep(5 * time.Millisecond)
				drawEvent, err = drawMio.Read()
				if err != nil {
					t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
				}
			}
			for _, move := range tc.p2Moves {
				t.Log("Input from P2 Mio")

				p2Mio.Write(io.InputEvent{Move: move, Terminate: false})
				// Read the Draw Output
				time.Sleep(5 * time.Millisecond)
				drawEvent, err = drawMio.Read()
				if err != nil {
					t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
				}
			}

			// Compare against the template
			if !strings.HasPrefix(drawEvent.DrawString, tc.expectedOutput) {
				t.Fatal("Moving Player 1 failed.")
			}
		})
	}
}

func TestUserMovement(t *testing.T) {

	// Test Wrapping Move Left
	testCases := []struct {
		name           string
		movements      []io.Move
		expectedOutput string
	}{
		{"Move User Left - Wrapping", []io.Move{io.Move_Left}, "-------------\r\n| . | . | . |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n| . | . | 1 |\r\n-------------\r\n"},
		{"Move User Down - Wrapping", []io.Move{io.Move_Left, io.Move_Down}, "-------------\r\n| . | . | 1 |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n"},
		{"Move User Right - Wrapping", []io.Move{io.Move_Left, io.Move_Down, io.Move_Right}, "-------------\r\n| . | . | . |\r\n-------------\r\n| 1 | . | . |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n"},
		{"Move User Up - Nonwrapping", []io.Move{io.Move_Left, io.Move_Down, io.Move_Right, io.Move_Up}, "-------------\r\n| 1 | . | . |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n"},
		{"Place Piece Move", []io.Move{io.Move_Left, io.Move_Down, io.Move_Right, io.Move_Up, io.Move_PlaceMark}, "-------------\r\n| X | . | . |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n| . | . | 2 |\r\n-------------\r\n"},
		// {"Player 2 Move Up", []io.Move{io.Move_Left, io.Move_Down, io.Move_Right, io.Move_Up, io.Move_PlaceMark, io.Move_Up}, "-------------\r\n| X | . | . |\r\n-------------\r\n| . | . | 2 |\r\n-------------\r\n| . | . | . |\r\n-------------\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Test
			p1Mio, _, drawMio, gameObject := SetupTest(t)
			defer gameObject.CloseGame()
			defer TeardownTest(t)

			// Discard the first read
			drawMio.Read()

			// Wait for users turn and perform movement
			var drawEvent io.DrawEvent
			var err error
			for _, move := range tc.movements {
				t.Log("Input from P1 Mio")
				p1Mio.Write(io.InputEvent{Move: move, Terminate: false})

				// Read the Draw Output
				time.Sleep(5 * time.Millisecond)
				drawEvent, err = drawMio.Read()
				if err != nil {
					t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
				}
			}

			// Compare against the template
			if !strings.HasPrefix(drawEvent.DrawString, tc.expectedOutput) {
				t.Fatal("Moving Player 1 failed.")
			}
		})
	}

}
func TestLocalGame(t *testing.T) {

	var GridTemplate = regexp.MustCompile(`-{13}(\r\n\|( [\d\.] \|){3}\r\n-{13}){3}`)
	_, _, drawMio, gameObject := SetupTest(t)
	defer gameObject.CloseGame()
	defer TeardownTest(t)

	// Read and compare against the template
	drawEvent, err := drawMio.Read()
	if err != nil {
		t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
	}
	templateMatch := GridTemplate.MatchString(drawEvent.DrawString)
	if !templateMatch {
		t.Fatalf("The draw event did not match the expected grid layout")
	}
}
