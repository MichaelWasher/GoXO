package main

import (
	"log"
	"os"
)

// TODO Add Multiplayer Support
// TODO Add Socket Support for Multiple Player Input
// TODO Flags for the Socket Connection


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
	// Configure Logging
	logFile := initLog()
	defer logFile.Close()

	// Configure Game
	game := Game{}
	game.InitGame()
	defer game.CloseGame()

	game.gameLoop()

}
