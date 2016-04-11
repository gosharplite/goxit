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

func (chain *Chain) Initialize(size int) {

	chain.SIZE = size
	chain.BOARD_SIZE = (chain.SIZE+2)*(chain.SIZE+1) + 1

	chain.MAX_POINTS = chain.SIZE * chain.SIZE
	chain.MAX_LIBERTIES = chain.SIZE * chain.SIZE

	chain.num_points = 0
	chain.points = make([]int, chain.MAX_POINTS)
	chain.points_indices = make([]int, chain.BOARD_SIZE)

	chain.num_liberties = 0
	chain.liberties = make([]int, chain.MAX_LIBERTIES)
	chain.liberties_indices = make([]int, chain.BOARD_SIZE)

	for i := 0; i < chain.BOARD_SIZE; i++ {

		chain.points_indices[i] = -1

		chain.liberties_indices[i] = -1
	}
}

func (chain *Chain) addPoint(point int) {

	// if point is in chain, do nothing
	if chain.points_indices[point] != -1 {
		return
	}

	chain.points[chain.num_points] = point
	chain.points_indices[point] = chain.num_points

	chain.num_points++
}

func (chain *Chain) addLiberty(point int) {

	// if point is in chain, do nothing
	if chain.liberties_indices[point] != -1 {
		return
	}

	chain.liberties[chain.num_liberties] = point

	chain.liberties_indices[point] = chain.num_liberties

	chain.num_liberties++
}

func (chain *Chain) hasPoint(point int) bool {

	return chain.points_indices[point] != -1
}

func (chain *Chain) removeLiberty(point int) {

	// if point is not in chain, do nothing
	if chain.liberties_indices[point] == -1 {
		return
	}

	// swap last liberty with current liberty
	index := chain.liberties_indices[point]
	end_liberty := chain.liberties[chain.num_liberties-1]
	chain.liberties[index] = end_liberty
	chain.liberties_indices[end_liberty] = index

	// remove point
	chain.liberties[chain.num_liberties-1] = 0
	chain.liberties_indices[point] = -1
	chain.num_liberties--
}
