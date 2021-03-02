package io

import "context"

type Move int32

const (
	Move_Noop      Move = 0
	Move_Left      Move = 1
	Move_Right     Move = 2
	Move_Up        Move = 3
	Move_Down      Move = 4
	Move_PlaceMark Move = 5
	Move_Quit      Move = 6
	Request_Move   Move = 7
)

type InputEvent struct {
	Move      Move
	Terminate bool
}
type DrawEvent struct {
	DrawString string
	Terminate  bool
}

func NewInputEvent(move Move) InputEvent {
	return InputEvent{Move: move, Terminate: false}
}

func NewDrawEvent(outputString string, terminate bool) DrawEvent {
	return DrawEvent{DrawString: outputString, Terminate: terminate}
}

type OutputHandler interface {
	RegisterDrawEvents(context.Context, <-chan DrawEvent)
}

type PlayerInputHandler interface {
	RegisterInputEvents(context.Context, chan InputEvent)
}
