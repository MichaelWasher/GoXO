package game_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/MichaelWasher/GoXO/pkg/game"
	"github.com/MichaelWasher/GoXO/pkg/io"
)

type MockGameIO struct {
	InputChannel chan io.InputEvent
	DrawChannel  <-chan io.DrawEvent
}

func (mio *MockGameIO) RegisterDrawEvents(ctx context.Context, drawChannel <-chan io.DrawEvent) {
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
func (mio MockGameIO) Write(ioe io.InputEvent) error {
	if mio.InputChannel == nil {
		return errors.New("Input Channel has not been configured correctly")
	}
	mio.InputChannel <- ioe
	return nil
}

func (mio MockGameIO) Read() (io.DrawEvent, error) {
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
func TestLocalGame(t *testing.T) {

	// Configure the IO
	p1Mio := &MockGameIO{}
	p2Mio := &MockGameIO{}

	// Create the Game
	log.Println("Game Created")
	gameObject := game.NewGame(p1Mio, p2Mio)
	t.Log("Game Created")
	defer gameObject.CloseGame()

	// Setup Game Loop
	go gameObject.GameLoop()
	for {
		// TODO Implement exponential Standoff
		time.Sleep(10 * time.Millisecond)
		if p1Mio.DrawChannel != nil {
			break
		}
	}
	drawEvent, err := p2Mio.Read()
	if err != nil {
		t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
	}
	// Compare against the template
	templateMatch := GridTemplate.MatchString(drawEvent.DrawString)
	if !templateMatch {
		t.Fatalf("The draw event did not match the expected grid layout")
	}

	// Test Basic Move Right
	expectedOutput := `-------------
| . | 1 | . |
-------------
| . | . | . |
-------------
| . | . | . |
-------------`
	p1Mio.Write(io.InputEvent{Move: io.Move_Right})

	time.Sleep(10 * time.Millisecond)
	drawEvent, err = p2Mio.Read()
	if err != nil {
		t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
	}
	// Compare against the template
	if strings.HasPrefix(drawEvent.DrawString, expectedOutput) {
		t.Fatal("Moving Player 1 failed.")
	}
	// Test Basic Move Left
	expectedOutput = `-------------
| 1 | . | . |
-------------
| . | . | . |
-------------
| . | . | . |
-------------`
	p1Mio.Write(io.InputEvent{Move: io.Move_Right})

	time.Sleep(10 * time.Millisecond)
	drawEvent, err = p2Mio.Read()
	if err != nil {
		t.Fatalf("Unable to read draw event. Received Error: %v. Exected event; Got %v", err, drawEvent)
	}
	// Compare against the template
	if strings.HasPrefix(drawEvent.DrawString, expectedOutput) {
		t.Fatal("Moving Player 1 failed.")
	}
}

//
var GridTemplate = regexp.MustCompile(`-{13}(\r\n\|( [\d\.] \|){3}\r\n-{13}){3}`)
