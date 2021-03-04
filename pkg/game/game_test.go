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
	// TODO - Replace with a channel merge and select
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

type command struct {
	player string
	move   io.Move
}
type gridSquare struct {
	index int
	value rune
}

func getBoard(gridSquares []gridSquare) string {
	const GRID_TEMPLATE = "" + // NOTE: Requires manual `\r\n` as Terminal will be in RAW mode and not pre-baked / cooked
		"-------------\r\n" +
		"| %s | %s | %s |\r\n" +
		"-------------\r\n" +
		"| %s | %s | %s |\r\n" +
		"-------------\r\n" +
		"| %s | %s | %s |\r\n" +
		"-------------\r\n"

	// Configure Replacements for the Placements
	replacementValues := make([]interface{}, 9)
	for i := range replacementValues {
		replacementValues[i] = "."
	}
	for _, gs := range gridSquares {
		replacementValues[gs.index] = string(gs.value)
	}

	// Replace and Return
	return fmt.Sprintf(GRID_TEMPLATE, replacementValues...)
}

func TestGeneric(t *testing.T) {
	// Test Wrapping Move Left
	var player1 = "p1"
	var player2 = "p2"

	testCases := []struct {
		name           string
		moves          []command
		expectedOutput string
	}{
		{"Check Game Template", []command{{player1, io.Move_Noop}}, getBoard([]gridSquare{{0, '1'}})},
		{"Player 2 Move Up", []command{
			{player: player1, move: io.Move_PlaceMark},
			{player: player2, move: io.Move_Up}},
			getBoard([]gridSquare{{0, 'X'}, {6, '2'}}),
		},
		{"Move User Left - Wrapping", []command{
			{player1, io.Move_Left}}, getBoard([]gridSquare{{8, '1'}})},
		{"Move User Down - Wrapping", []command{
			{player1, io.Move_Left}, {player1, io.Move_Down}},
			getBoard([]gridSquare{{2, '1'}})},
		{"Move User Right - Wrapping", []command{
			{player1, io.Move_Left}, {player1, io.Move_Down}, {player1, io.Move_Right}},
			getBoard([]gridSquare{{3, '1'}})},
		{"Move User Up - Nonwrapping", []command{
			{player1, io.Move_Left}, {player1, io.Move_Down}, {player1, io.Move_Right}, {player1, io.Move_Up}},
			getBoard([]gridSquare{{0, '1'}})},
		{"Place Piece Move", []command{
			{player1, io.Move_Left}, {player1, io.Move_Down}, {player1, io.Move_Right}, {player1, io.Move_PlaceMark}},
			getBoard([]gridSquare{{0, '2'}, {3, 'X'}})},
		{"Muti-player - Player Alternating ", []command{
			{player1, io.Move_Left}, {player1, io.Move_Down}, {player1, io.Move_Right}, {player1, io.Move_Up}, {player1, io.Move_PlaceMark}, {player2, io.Move_Up}},
			getBoard([]gridSquare{{0, 'X'}, {6, '2'}})},
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
			for _, command := range tc.moves {
				// Configure the current command IO
				var currentMio *MockGameIO
				if command.player == player1 {
					currentMio = p1Mio
				} else {
					currentMio = p2Mio
				}

				t.Log("Input from P1 Mio")
				currentMio.Write(io.InputEvent{Move: command.move, Terminate: false})

				// Read the Draw Output
				time.Sleep(5 * time.Millisecond)
				drawEvent, err = drawMio.Read()
				if err != nil {
					t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
				}
			}

			// Compare against the template
			if !strings.HasPrefix(drawEvent.DrawString, tc.expectedOutput) {
				t.Fatalf("%v has failed. Expected board:\n%v\nGot:\n%v\n", tc.name, tc.expectedOutput, drawEvent.DrawString)
			}
		})
	}
}

// TODO Test Perform Win
// TODO Create Generic Test that can be used for iterative moves

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
