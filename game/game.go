package game

import (
	"fmt"
	"github.com/pkg/term"
	"log"
	"math"
)

// -- Use the Protobuf values for Move Struct


// ---- Game Variables
var lg LogicGrid = NewLogicGrid()

// Users
var player1 = User{Position: 0, Character: "Y", Mark: "X", Name: "Player1"}
var player2 = User{Position: 8, Character: "Y", Mark: "0", Name: "Player2"}
var players = []*User{&player1, &player2}

var currentPlayerIndex int = 0


type Move int32

const (
	//Move_Noop      Move = 0
	Move_Left      Move = 1
	Move_Right     Move = 2
	Move_Up        Move = 3
	Move_Down      Move = 4
	Move_PlaceMark Move = 5
	Move_Quit      Move = 6
)


type Game struct {
	Terminal *term.Term
	Running  bool
}

func (game *Game) InitGame() {
	game.Terminal, _ = term.Open("/dev/tty")

	// Configure Terminal
	term.RawMode(game.Terminal)


}
func (game *Game) CloseGame() {
	defer game.Terminal.Close() // Defer is LIFO ordering, Close is last.
	defer game.Terminal.Restore()
}
// Core Game Loop
var OutstandingMoves = make(chan Move)

func (game *Game) GameLoop() {
	game.Running = true
	game.draw()

	//// Setup Game input
	for game.Running {
		currentMove := <- OutstandingMoves
		// Iterate Game States
		game.update(currentMove)
		game.draw()
	}
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

	currentPlayer := players[currentPlayerIndex]
	if currentMove == Move_Quit{
		game.Running = false
	}else if currentMove == Move_PlaceMark{
		lg.PlaceMark(currentPlayer)
		if game.checkWinner(currentPlayer){
			//TODO Break and start closing sequence

		}
		// Change Player
		currentPlayerIndex = (currentPlayerIndex + 1) % len(players)
	}else{
		performMove(lg, currentPlayer, currentMove)
	}

}

func (game *Game) draw() {
	lg.draw([]*User{players[currentPlayerIndex]})
	//-- Draw Statistics
	statsTemplate := "Player Name: %s\r\nCurrent Position: %d\r\n"
	for _, player := range players{
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
		for _, v := range column{
			log.Print("%d,",v)
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
