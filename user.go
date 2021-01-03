package main

import "math"

// User represents a single player within the game grid
type User struct {
	Position  int
	Character string
	Mark      string
	Name string
}



// MoveUser will change the position of the user, if possible, in the provided Move (d)
func (player *User) MoveUser(d Move) {
	switch d {
	case MoveLeft:
		if player.Position%3 != 0 {
			player.Position--
		}
	case MoveRight:
		if player.Position%3 != 2 {
			player.Position++
		}
	case MoveUp:
		if math.Floor(float64(player.Position/3)) != 0 {
			player.Position -= 3
		}
	case MoveDown:
		if math.Floor(float64(player.Position/3)) != 2 {
			player.Position += 3
		}
	default:
		println("Error has occurred in the move user function.")
	}
}
