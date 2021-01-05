package main

import (
	"github.com/pkg/term"
	"log"
	"os"
)

// TODO Add Multiplayer Support
// TODO Add Socket Support for Multiple Player Input
// TODO Flags for the Socket Connection

var terminal, _ = term.Open("/dev/tty")


// Configure Loggin
func initLog() *os.File{
	f, err := os.OpenFile("log-file.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.Println("This is a test log entry")
	return f
}

func main() {
	// Configure Logs
	logFile := initLog()
	defer logFile.Close()
	// Configure Terminal
	term.RawMode(terminal)
	defer terminal.Close() // Defer is LIFO ordering, Close is last.
	defer terminal.Restore()


	go setup_server(7777)
	gameLoop()
}
