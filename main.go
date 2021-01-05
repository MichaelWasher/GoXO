package main


// TODO Add Multiplayer Support
// TODO Add Socket Support for Multiple Player Input
// TODO Flags for the Socket Connection


func main() {
	// Configure Logs
	game := Game{}
	game.InitGame()
	defer game.CloseGame()

	game.gameLoop()
}
