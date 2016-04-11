package board

import ()

type MoveHistory struct {
	player BoardState
	point  int

	// Ko point before move was played
	ko_point int

	// capture directions[d] = true if and only if
	// a capture occurred in the direction d from point
	capture_directions []bool
}

func (history *MoveHistory) Initialize(player BoardState, point int, ko_point int) {

	history.player = player

	history.point = point

	history.ko_point = ko_point

	history.capture_directions = []bool{false, false, false, false}
}

func (history *MoveHistory) setCapture_directions(dir int) {

	history.capture_directions[dir] = true
}

func (history *MoveHistory) isCapture_directions(dir int) bool {

	return history.capture_directions[dir]
}
