package game

import (
	"fmt"
	"github.com/pkg/term"
	"log"
	"math"
	// TODO Work out the dependency into an interface
	. "github.com/MichaelWasher/GoXO/grpc"
)

// -- Use the Protobuf values for Move Struct


// ---- Game Variables
var lg LogicGrid = NewLogicGrid()

var currentPlayerIndex int = 0

type Game struct {
	Terminal *term.Term
	Running  bool
	// Users
	player1 User
	player2 User
	players []*User
	// Display Other Users
	displayOtherUsers bool
}

func (game *Game) InitGame() {
	game.Terminal, _ = term.Open("/dev/tty")

	// Configure Terminal
	term.RawMode(game.Terminal)

	// Configure Defaults
	game.displayOtherUsers = true

	// Configure Player Inputs
	playerOneInput := make(chan Move)
	playerTwoInput := make(chan Move)

	// Init Players
	game.player1 = User{Position: 0, Character: "1", Mark: "X", Name: "Player1", InputChannel: &playerOneInput}
	game.player2 = User{Position: 8, Character: "2", Mark: "0", Name: "Player2", InputChannel: &playerTwoInput}
	game.players = []*User{&game.player1, &game.player2}
}
func (game *Game) CloseGame() {
	defer game.Terminal.Close() // Defer is LIFO ordering, Close is last.
	defer game.Terminal.Restore()
}
// Core Game Loop
func (game *Game) GameLoop() {

	game.Running = true
	game.draw()
	var currentMove Move

	// Check for input from either player
	for game.Running {
		select{
		case currentMove = <- *game.player1.InputChannel:
			log.Println("Received input from Player 1")
		case currentMove = <- *game.player2.InputChannel:
			log.Println("Received input from Player 2")
		}

		// Iterate Game States
		game.update(currentMove)
		game.draw()
	}
}

func (game *Game) GetPlayerOneInputChannel() *chan Move {
	return game.player1.InputChannel
}

func (game *Game) GetPlayerTwoInputChannel() *chan Move {
	return game.player2.InputChannel
}


func performMove(lg LogicGrid, currentPlayer *User, currentMove Move){
	tmpPosition := currentPlayer.Position
	log.Print("PerformMove")

	// TODO Fix this to deal with corner's better. If top left and bottom right are only left, cannot jump between.
	for{
		if (tmpPosition%3 == 0 && Move_Left == currentMove) ||
			(tmpPosition%3 == 2 && Move_Right == currentMove) ||
			(math.Floor(float64(tmpPosition/3)) == 0 && Move_Up == currentMove) ||
			(math.Floor(float64(tmpPosition/3)) == 2 && Move_Down == currentMove){
			return
		}

		switch currentMove {
		case Move_Left:
				tmpPosition--
		case Move_Right:
				tmpPosition++
		case Move_Up:
				tmpPosition -= 3
		case Move_Down:
				tmpPosition += 3
		default:
			println("Error has occurred in the move user function.")
		}

		if !lg.isUsed(Position(tmpPosition)){ //TODO Replace position type
			currentPlayer.Position = tmpPosition
			log.Printf("Current Player Position %d", tmpPosition)
			return
		}
	}

}

func (game *Game) update(currentMove Move) {

	currentPlayer := game.players[currentPlayerIndex]
	if currentMove == Move_Quit{
		game.Running = false
	}else if currentMove == Move_PlaceMark{
		lg.PlaceMark(currentPlayer)
		if game.checkWinner(currentPlayer){
			//TODO Break and start closing sequence

		}
		// Change Player
		currentPlayerIndex = (currentPlayerIndex + 1) % len(game.players)
	}else{
		performMove(lg, currentPlayer, currentMove)
	}

}

func (game *Game) draw() {
	if game.displayOtherUsers {
		lg.draw([]*User{game.players[currentPlayerIndex]})
	} else {
		lg.draw(game.players)
	}

	//-- Draw Statistics
	statsTemplate := "Player Name: %s\r\nCurrent Position: %d\r\n"
	for _, player := range game.players{
		fmt.Printf(statsTemplate, player.Name, player.Position)
	}
}


// ---- Check Winner Functionality
func (game *Game) checkWinner(user *User)(bool) {
	if columnsComplete(lg) || rowsComplete(lg) || diagComplete(lg) || antiDiagComplete(lg){
		// Winner is found
		log.Print("WINNER WAS FOUND")
		return true
	}
	return false
}

func rowsComplete(lg LogicGrid) bool{
	for currentRow := 0; currentRow < RowLength; currentRow++ {
		row, _ := lg.getRow(currentRow)
		if checkAllSame(row){
			return true
		}
	}
	return false
}
func columnsComplete(lg LogicGrid) bool{
	for currentColumn := 0; currentColumn < RowCount; currentColumn++ {
		column, _ := lg.getColumn(currentColumn)
		log.Printf("Checking Columns")
		for i, v := range column{
			log.Printf("%d,%s", i, v)
		}
		if checkAllSame(column){
			return true
		}
	}
	log.Printf("Colmns not complete")
	return false
}

func diagComplete (lg LogicGrid) bool{
	diagonal := lg.getDiagonal()
	return checkAllSame(diagonal)
}

func antiDiagComplete (lg LogicGrid) bool{
	diagonal := lg.getAntiDiagonal()
	return checkAllSame(diagonal)
}

// -- helper function
func checkAllSame(a []string) bool {
	for i := 1; i < len(a); i++ {
		if a[i] != a[0] {
			return false
		}
	}
	return true
}
