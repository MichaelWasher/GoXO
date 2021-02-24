package game

import (
	"context"
	"fmt"
	"log"

	"github.com/MichaelWasher/GoXO/pkg/io"
)

// ---- Game Variables
var lg LogicGrid = NewLogicGrid()

var currentPlayerIndex int = 0

type State int

const (
	StartGame State = iota
	MidGame
	EndGame
)

type Game struct {
	Running bool
	State   State
	// Users
	player1       User
	player2       User
	players       []*User
	currentPlayer *User
	// Display Other Users
	drawChannel       chan io.DrawEvent
	displayOtherUsers bool
	// IO Channels
	// Context Management
	context context.Context
	cancel  context.CancelFunc
}

func NewGame(player1 io.PlayerEventHandler, player2 io.PlayerEventHandler) *Game {
	game := Game{}
	// Configure Defaults
	game.displayOtherUsers = true
	game.context, game.cancel = context.WithCancel(context.Background())

	// Configure Player Outputs
	game.drawChannel = make(chan io.DrawEvent, 1)
	go player1.RegisterDrawEvents(game.context, game.drawChannel)
	go player2.RegisterDrawEvents(game.context, game.drawChannel)

	// Configure Player Inputs
	player1Input := make(chan io.InputEvent, 1)
	go player1.RegisterInputEvents(game.context, player1Input)

	player2Input := make(chan io.InputEvent, 1)
	go player2.RegisterInputEvents(game.context, player2Input)

	// TODO Allow for Multiple Event Subscribers

	// Init Players
	game.player1 = User{Position: 0, Character: "1", Mark: "X", Name: "Player1", PlayerEventHandler: player1, DrawChannel: game.drawChannel, InputChannel: player1Input}
	game.player2 = User{Position: 8, Character: "2", Mark: "0", Name: "Player2", PlayerEventHandler: player2, DrawChannel: game.drawChannel, InputChannel: player2Input}

	game.players = []*User{&game.player1, &game.player2}

	return &game
}

func (game *Game) CloseGame() {
	game.cancel()
}

// Core Game Loop
func (game *Game) GameLoop() {

	game.Running = true
	game.draw()

	// Check for input from either player
	for game.Running {

		// Iterate Game States
		game.update()
		game.draw()
	}
	game.CloseGame()
}

func performMove(lg LogicGrid, currentPlayer *User, currentMove io.Move) {
	tmpPosition := currentPlayer.Position
	log.Print("PerformMove")
	// TODO Fix this to deal with corner's better. If top left and bottom right are only left, cannot jump between.
	for {

		switch currentMove {
		case io.Move_Left:
			tmpPosition--
		case io.Move_Right:
			tmpPosition++
		case io.Move_Up:
			tmpPosition -= 3
		case io.Move_Down:
			tmpPosition += 3
		default:
			println("Error has occurred in the move user function.")
		}
		// Add Wrapping
		if tmpPosition < 0 {
			tmpPosition += (RowCount * RowLength)
		}
		if tmpPosition >= (RowCount * RowLength) {
			tmpPosition -= (RowCount * RowLength)
		}

		// Skip Used Squares
		if !lg.isUsed(Position(tmpPosition)) {
			currentPlayer.Position = tmpPosition
			log.Printf("Current Player Position %d", tmpPosition)
			return
		}
	}

}

func (game *Game) update() {
	// Player Turn
	currentPlayer := game.players[currentPlayerIndex]
	currentPlayer.InputChannel <- io.NewInputEvent(io.Request_Move)
	// Perform Turn
	event := <-currentPlayer.InputChannel
	log.Printf("Received input from %s", currentPlayer.Name)

	switch event.Move {
	case io.Move_Quit:
		log.Print("Command Received to End Game.")
		game.Running = false
	case io.Move_PlaceMark:
		lg.PlaceMark(currentPlayer)
		if game.checkWinner(currentPlayer) {
			//TODO Break and start closing sequence

		}
		// Change Player
		currentPlayerIndex = (currentPlayerIndex + 1) % len(game.players)
	default:
		performMove(lg, currentPlayer, event.Move)
	}

}

func (game *Game) draw() {
	var outputString string
	if game.displayOtherUsers {
		outputString = lg.draw([]*User{game.players[currentPlayerIndex]})
	} else {
		outputString = lg.draw(game.players)
	}

	//-- Draw Statistics
	statsTemplate := "Player Name: %s\r\nCurrent Position: %d\r\n"
	for _, player := range game.players {
		outputString += fmt.Sprintf(statsTemplate, player.Name, player.Position)
	}
	// Perform a Draw Event through the Draw Channel
	game.drawChannel <- io.NewDrawEvent(outputString, !game.Running)
}

// ---- Check Winner Functionality
func (game *Game) checkWinner(user *User) bool {
	if columnsComplete(lg) || rowsComplete(lg) || diagComplete(lg) || antiDiagComplete(lg) {
		// Winner is found
		log.Print("WINNER WAS FOUND")
		return true
	}
	return false
}

func rowsComplete(lg LogicGrid) bool {
	for currentRow := 0; currentRow < RowLength; currentRow++ {
		row, _ := lg.getRow(currentRow)
		if checkAllSame(row) {
			return true
		}
	}
	return false
}
func columnsComplete(lg LogicGrid) bool {
	for currentColumn := 0; currentColumn < RowCount; currentColumn++ {
		column, _ := lg.getColumn(currentColumn)
		log.Printf("Checking Columns")
		for i, v := range column {
			log.Printf("%d,%s", i, v)
		}
		if checkAllSame(column) {
			return true
		}
	}
	log.Printf("Colmns not complete")
	return false
}

func diagComplete(lg LogicGrid) bool {
	diagonal := lg.getDiagonal()
	return checkAllSame(diagonal)
}

func antiDiagComplete(lg LogicGrid) bool {
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
