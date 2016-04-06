package board

import ()

type Chain struct {
	SIZE       int //Size of Go board.
	BOARD_SIZE int //(SIZE+2)*(SIZE+1)+1

	// SIZE*SIZE is very loose upper bound on the number of
	// points and liberties that a chain can have
	MAX_POINTS    int //SIZE*SIZE
	MAX_LIBERTIES int //SIZE*SIZE

	// Data members for keeping track of points
	points         []int //MAX_POINTS
	num_points     int
	points_indices []int //BOARD_SIZE;

	// Data members for keeping track of liberties
	liberties         []int //MAX_LIBERTIES
	num_liberties     int
	liberties_indices []int //BOARD_SIZE
}
