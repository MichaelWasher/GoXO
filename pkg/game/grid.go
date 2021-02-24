package game

import (
	"errors"
	"fmt"
	"log"
	"math"
)

// --- Constants
const GRID_TEMPLATE = "" + // NOTE: Requires manual `\r\n` as Terminal will be in RAW mode and not pre-baked / cooked
	"-------------\r\n" +
	"| %s | %s | %s |\r\n" +
	"-------------\r\n" +
	"| %s | %s | %s |\r\n" +
	"-------------\r\n" +
	"| %s | %s | %s |\r\n" +
	"-------------\r\n"

// ---- Structures
const RowCount int = 3
const RowLength int = 3

type LogicGrid struct {
	GridArray    [RowCount * RowLength]string
	defaultValue string
}
type Position int

// --- Constructor(s)
func NewLogicGrid() LogicGrid {
	var lg LogicGrid
	lg.defaultValue = "."
	for i := 0; i < len(lg.GridArray); i++ {
		lg.GridArray[i] = lg.defaultValue
	}
	return lg
}

// --- Methods

// Check if position in grid is used
func (lg LogicGrid) isUsed(currentPosition Position) bool {
	if lg.GridArray[currentPosition] == lg.defaultValue {
		return false
	}
	return true
}

// Get the Display Grid used for Drawing
func (lg LogicGrid) getDisplayGrid() []interface{} {
	displayGrid := make([]interface{}, RowCount*RowLength)
	for i, v := range lg.GridArray {
		displayGrid[i] = v
	}
	return displayGrid
}
func (lg LogicGrid) getRow(index int) ([]string, error) {
	if index > RowCount-1 {
		log.Fatal("Index out of bounds")
		return nil, errors.New("Index out of bounds")
	}

	startIndex, endIndex := index*RowCount, index*RowCount+RowCount
	return lg.GridArray[startIndex:endIndex], nil
}

func (lg LogicGrid) getColumn(colIndex int) ([]string, error) {
	if colIndex > RowLength-1 {
		log.Fatal("Index out of bounds")
		return nil, errors.New("Index out of bounds")
	}

	column := make([]string, RowCount)
	for index := range column {
		column[index] = lg.GridArray[index*RowLength+colIndex]
	}

	return column, nil
}

func (lg LogicGrid) getDiagonal() []string {
	diagonal := make([]string, int(math.Min(float64(RowLength), float64(RowCount))))
	for index := range diagonal {
		diagonal[index] = lg.GridArray[(index*RowLength)+index]
	}

	return diagonal
}

func (lg LogicGrid) getAntiDiagonal() []string {
	diagonal := make([]string, int(math.Min(float64(RowLength), float64(RowCount))))
	for index := range diagonal {
		diagonal[index] = lg.GridArray[(RowLength-index-1)+(RowLength*index)]
	}
	return diagonal
}

// PlaceMark places one of the users tokens on the currently occupied square
func (lg *LogicGrid) PlaceMark(player *User) error {
	// Ensure that the space is not used
	if lg.isUsed(Position(player.Position)) {
		return errors.New("Invalid Position. Unable to place mark.")
	}
	log.Print("Placing Mark")
	log.Printf("Mark is %s", lg.GridArray[player.Position])
	lg.GridArray[player.Position] = player.Mark
	log.Printf("Mark is %s", lg.GridArray[player.Position])
	return nil
}

// Return a string of the display screen
func (lg LogicGrid) draw(players []*User) string {

	// Get Grid for Display
	displayGrid := lg.getDisplayGrid()

	// Add Players
	for _, player := range players {
		displayGrid[player.Position] = player.Character
	}

	// Print all Elements
	return fmt.Sprintf(GRID_TEMPLATE, displayGrid...)
}
