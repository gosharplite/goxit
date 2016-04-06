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
