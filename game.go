package main

import (
	"log"
	"math"
	"time"
)

type Move int
const (
	MoveLeft Move = 1 << iota
	MoveRight
	MoveUp
	MoveDown
	PlacePiece
	Quit
)


// ---- Game Variables

var lg LogicGrid = NewLogicGrid()

// Users
var player1 = User{Position: 0, Character: "Y", Mark: "X"}
var player2 = User{Position: 8, Character: "Y", Mark: "0"}
var players = []*User{&player1, &player2}
var currentUser = players[0]

var running bool

// Move State Carriers
var lastMove Move


// Core Game Loop
var outstandingMoves = make(chan Move)
func gameLoop() {
	running = true
	go handleKeyEvents()
	for running {
		draw()
		update()
		time.Sleep(1)
	}
}

func performMove(lg LogicGrid, currentPlayer *User, currentMove Move){
	tmpPosition := currentPlayer.Position
	log.Print("PerformMove")

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

func update() {
	currentMove := <- outstandingMoves
	currentPlayer := players[0]
	if currentMove == Quit{
		running = false
	}else if currentMove == PlacePiece{
		// Place Mark in Grid and cycle current user

	}else{
		performMove(lg, currentPlayer, currentMove)
	}

}

func draw() {
	lg.draw(players)
}


// TODO -- Implement the game logic. Wining, Losing and Points
func  checkWinner(user User)(bool) {
	if columnsComplete(lg) || rowsComplete(lg) || diagComplete(lg) || antiDiagComplete(lg){
		// Winnner is found
	}
	return false
}

// ---- Utility Functions

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
		if checkAllSame(column){
			return true
		}
	}
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
