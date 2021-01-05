package main

import (
	tm "github.com/buger/goterm"
	"log"
	"math"
)

type Move int
const (
	MoveLeft Move = 1 << iota
	MoveRight
	MoveUp
	MoveDown
	PlacePiece
	Quit
	FinishedMove
)


// ---- Game Variables
var lg LogicGrid = NewLogicGrid()

// Users
var player1 = User{Position: 0, Character: "Y", Mark: "X", Name: "Player1"}
var player2 = User{Position: 8, Character: "Y", Mark: "0", Name: "Player2"}
var players = []*User{&player1, &player2}

// Core Game Loop
var outstandingMoves = make(chan Move)
var currentPlayerIndex int = 0
var running bool

func gameLoop() {
	running = true
	// Initial Draw
	draw()

	go handleKeyEvents()
	for running {
		currentMove := <- outstandingMoves
		update(currentMove)

		draw()
		outstandingMoves <- FinishedMove
	}
}

func performMove(lg LogicGrid, currentPlayer *User, currentMove Move){
	tmpPosition := currentPlayer.Position
	log.Print("PerformMove")

	// TODO Fix this to deal with corner's better. If top left and bottom right are only left, cannot jump between.
	for{
		if (tmpPosition%3 == 0 && MoveLeft == currentMove) ||
			(tmpPosition%3 == 2 && MoveRight == currentMove) ||
			(math.Floor(float64(tmpPosition/3)) == 0 && MoveUp == currentMove) ||
			(math.Floor(float64(tmpPosition/3)) == 2 && MoveDown == currentMove){
			return
		}

		switch currentMove {
		case MoveLeft:
				tmpPosition--
		case MoveRight:
				tmpPosition++
		case MoveUp:
				tmpPosition -= 3
		case MoveDown:
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

func update(currentMove Move) {

	currentPlayer := players[currentPlayerIndex]
	if currentMove == Quit{
		running = false
	}else if currentMove == PlacePiece{
		lg.PlaceMark(currentPlayer)
		if checkWinner(currentPlayer){
			//TODO Break and start closing sequence

		}
		// Change Player
		currentPlayerIndex = (currentPlayerIndex + 1) % len(players)
	}else{
		performMove(lg, currentPlayer, currentMove)
	}

}

func draw() {
	lg.draw([]*User{players[currentPlayerIndex]})
	//-- Draw Statistics
	statsTemplate := `
Player Name: %s
Current Position: %d`
	for _, player := range players{
		tm.Printf(statsTemplate, player.Name, player.Position)
	}
	tm.Flush()
}


// ---- Check Winner Functionality
func  checkWinner(user *User)(bool) {
	if columnsComplete(lg) || rowsComplete(lg) || diagComplete(lg) || antiDiagComplete(lg){
		// Winnner is found
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
	log.Printf("Columns not complete")
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
