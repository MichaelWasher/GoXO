package main

import (
	"errors"
	tm "github.com/buger/goterm"
	"log"
	"math"
)

// --- Constants
const GRID_TEMPLATE = `
-------------
| %s | %s | %s |
-------------
| %s | %s | %s |
-------------
| %s | %s | %s |
-------------
`

// ---- Structures
const RowCount int = 3
const RowLength int = 3

type LogicGrid struct {
	GridArray [RowCount * RowLength] string
	defaultValue string
}
type Position int

// --- Constructor(s)
func NewLogicGrid () LogicGrid{
	var lg LogicGrid
	lg.defaultValue = "."
	for i := 0; i < len(lg.GridArray); i++ {
		lg.GridArray[i] = lg.defaultValue
	}
	return lg
}

// --- Methods

// Check if position in grid is used
func (lg LogicGrid) isUsed (currentPosition Position) bool{
	if lg.GridArray[currentPosition] == lg.defaultValue{
		return false
	}
	return true
}
// Get the Display Grid used for Drawing
func (lg LogicGrid) getDisplayGrid() []interface{}{
	displayGrid := make([]interface{}, RowCount*RowLength)
	for i, v := range lg.GridArray {
		displayGrid[i] = v
	}
	return displayGrid
}
func (lg LogicGrid) getRow(index int) ([]string, error){
	if index > RowCount-1{
		log.Fatal("Index out of bounds")
		return nil, errors.New("Index out of bounds")
	}

	startIndex, endIndex := index *RowCount, index *RowCount+RowCount
	return lg.GridArray[startIndex:endIndex], nil
}

func (lg LogicGrid) getColumn(colIndex int) ([]string, error){
	if colIndex > RowLength-1{
		log.Fatal("Index out of bounds")
		return nil, errors.New("Index out of bounds")
	}

	column := make([]string, RowCount)
	for index := range column{
		column[index] = lg.GridArray[index * RowLength + index]
	}

	return column, nil
}

func (lg LogicGrid) getDiagonal() ([]string){
	diagonal := make([]string, int(math.Min(float64(RowLength), float64(RowCount))))
	for index := range diagonal{
		diagonal[index] = lg.GridArray[(index * RowLength) + index]
	}

	return diagonal
}

func (lg LogicGrid) getAntiDiagonal() ([]string){
	diagonal := make([]string, int(math.Min(float64(RowLength), float64(RowCount))))
	for index := range diagonal{
		diagonal[index] = lg.GridArray[(RowLength - index - 1) - (RowLength * index)]
	}
	return diagonal
}


// Draw all elements on to the Screen
func (lg LogicGrid) draw(players []*User) {
	// Clear current screen
	tm.Clear()
	tm.MoveCursor(0, 0)

	// Get Grid for Display
	displayGrid := lg.getDisplayGrid()

	// Add Players
	for _, player := range players{
		displayGrid[player.Position] = player.Character
	}

	// Print all Elements
	tm.Printf(GRID_TEMPLATE, displayGrid...)
	tm.Flush()
}
