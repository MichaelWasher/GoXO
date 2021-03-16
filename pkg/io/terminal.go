package io

// TODO Add Thread Safety
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

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
	PlaceMark: []byte{32},  //TODO Space
	Quit:      []byte{3},   //CTRL + C
}

type Tty struct {
	*term.Term
	Lock               *sync.Mutex
	InputChannel       <-chan []byte
	SubscriberChannels []chan<- []byte
	Context            context.Context
	CancelFunction     func()
}

var ttyLock = &sync.Mutex{}
var ttyStaticInstance *Tty

// Terminal Definition - Matches the IO Handler Interface
type Terminal struct {
	*Tty
	keymap KeyMap
}

func ensureTtySingleton() (*Tty, error) {
	ttyLock.Lock()
	defer ttyLock.Unlock()

	if ttyStaticInstance == nil {
		// Create Tty
		t, err := term.Open("/dev/ttys005")
		if err != nil {
			log.Fatalf("Unable to open a terminal. %v", err)
		}

		// Configure Terminal into Raw Mode
		err = t.SetRaw()

		if err != nil {
			log.Fatalf("Unable to set Terminal into Raw Mode. %v", err)
			log.Fatal("Attempting to continue but the game may not display correctly.")
		}
		err = t.SetReadTimeout(500 * time.Millisecond)
		if err != nil {
			log.Fatal("Unable to configure the Terminal read timeout.", err)
			return nil, errors.New("unable to set terminal read timeout")
		}
		ctx, cancel := context.WithCancel(context.TODO())

		ttyStaticInstance = &Tty{
			Term:               t,
			Lock:               ttyLock,
			SubscriberChannels: []chan<- []byte{},
			Context:            ctx,
			CancelFunction:     cancel,
		}
		// Start the fan-out channel reading
		ConfigureTtyChannel(ttyStaticInstance, ctx)
	}

	return ttyStaticInstance, nil
}

func NewTerminal(km KeyMap) (*Terminal, error) {
	// Ensure Tty is configured
	ttyInstance, err := ensureTtySingleton()
	if err != nil {
		log.Fatalf("Unable to configure tty. %v", err)
		return nil, err
	}
	// Create User-Terminal
	terminal := &Terminal{}
	terminal.Tty = ttyInstance
	terminal.keymap = km

	return terminal, nil
}

func (t Terminal) Close() {
	// Set terminal back to sane mode
	defer t.Term.Close()
	t.SetCbreak()
}

func (t Terminal) Print(outString string) {
	// t.Flush()
	t.Write([]byte(outString))
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
			// TODO This does not respect the current Tty access. Should print to the Tty but this has buffering issues.
			t.Print("\033[H\033[2J")
			// fmt.Print(event.DrawString)
			t.Print(event.DrawString)
		}

	}
}

func (t *Terminal) RegisterInputEvents(ctx context.Context, playerInput chan InputEvent) {
	// Add PlayerInput to the terminal subscribers

	characterInputChannel := AddTtySubscriber()
	for {
		select {
		case char := <-characterInputChannel:
			// TODO If Character is expected from keymap
			log.Printf("Register Input Events in : %d", goid())
			playerInput <- inputEventFromBytes(char, t.keymap)
		case <-ctx.Done():
			return
		}
	}
	// TODO add remove Tty subscriber

}

func inputEventFromBytes(c []byte, km KeyMap) InputEvent {

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

func AddTtySubscriber() <-chan []byte {
	//TODO Add removal of subscriber
	// Synchronise
	ttyLock.Lock()
	defer ttyLock.Unlock()
	// Add inputChannel to subscribers
	characterInputChannel := make(chan []byte, 3)
	ttyStaticInstance.SubscriberChannels = append(ttyStaticInstance.SubscriberChannels, characterInputChannel)
	return characterInputChannel
}
func ConfigureTtyChannel(tty *Tty, ctx context.Context) {
	go func() {
		for {

			select {
			case <-ctx.Done():

				return

			default:
				bytes := make([]byte, 3)
				numRead, err := tty.Read(bytes)

				// Did not read anything from terminal in time (No need to log as this is expected)
				if err == io.EOF {
					continue
				}

				if err != nil {
					log.Printf("Unable to read from the Terminal. %v", err)
					continue
				}

				// For each subscriber fan out the input
				log.Printf("ConfigureTtyChannel : %d", goid())
				ttyLock.Lock()
				for _, inputChannel := range tty.SubscriberChannels {
					log.Printf("ConfigureTtyChannel: Write to SubscriberChannel %v", inputChannel)
					inputChannel <- bytes[0:numRead]
				}
				ttyLock.Unlock()
			}
		}

	}()
}

func goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
