package main

import "math"

// User represents a single player within the game grid
type User struct {
	Position  int
	Character string
	Mark      string
}

// PlaceMark places one of the users tokens on the currently occupied square
func (player User) PlaceMark() {
	if originGrid[player.Position] != "." {
		return
	}

	originGrid[player.Position] = player.Mark
}

// MoveUser will change the position of the user, if possible, in the provided Direction (d)
func (player *User) MoveUser(d Direction) {
	switch d {
	case Left:
		if player.Position%3 != 0 {
			player.Position--
		}
	case Right:
		if player.Position%3 != 2 {
			player.Position++
		}
	case Up:
		if math.Floor(float64(player.Position/3)) != 0 {
			player.Position -= 3
		}
	case Down:
		if math.Floor(float64(player.Position/3)) != 2 {
			player.Position += 3
		}
	default:
		println("Error has occurred in the move user function.")
	}
}
