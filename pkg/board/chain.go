package board

import ()

// A chain represents a group of points.
type chain struct {

	// Size is used for array length estimation.
	size      int
	boardSize int

	// size*size is very loose upper bound on the number of
	// points and liberties that a chain can have
	maxPoints    int
	maxLiberties int

	// Data members for keeping track of points
	// points length = maxPoints
	// pointsIndices length = boardSize
	points        []int
	numPoints     int
	pointsIndices []int

	// Data members for keeping track of liberties
	// liberties length = maxLiberties
	// libertiesIndices length = boardSize
	liberties        []int
	numLiberties     int
	libertiesIndices []int
}

func newChain(size int) chain {

	c := chain{}
	c.init(size)

	return c
}

func (c *chain) init(size int) {

	c.size = size
	c.boardSize = (c.size+2)*(c.size+1) + 1

	c.maxPoints = c.size * c.size
	c.maxLiberties = c.size * c.size

	c.numPoints = 0
	c.points = make([]int, c.maxPoints)
	c.pointsIndices = make([]int, c.boardSize)

	c.numLiberties = 0
	c.liberties = make([]int, c.maxLiberties)
	c.libertiesIndices = make([]int, c.boardSize)

	for i := 0; i < c.boardSize; i++ {

		c.pointsIndices[i] = -1

		c.libertiesIndices[i] = -1
	}
}

func (c *chain) addPoint(pt int) {

	// if point is in chain, do nothing
	if c.pointsIndices[pt] != -1 {
		return
	}

	c.points[c.numPoints] = pt
	c.pointsIndices[pt] = c.numPoints

	c.numPoints++
}

func (c *chain) addLiberty(pt int) {

	// if point is in chain, do nothing
	if c.libertiesIndices[pt] != -1 {
		return
	}

	c.liberties[c.numLiberties] = pt

	c.libertiesIndices[pt] = c.numLiberties

	c.numLiberties++
}

func (c *chain) hasPoint(pt int) bool {

	return c.pointsIndices[pt] != -1
}

func (c *chain) removeLiberty(pt int) {

	// if point is not in chain, do nothing
	if c.libertiesIndices[pt] == -1 {
		return
	}

	// swap last liberty with current liberty
	i := c.libertiesIndices[pt]
	j := c.liberties[c.numLiberties-1]
	c.liberties[i] = j
	c.libertiesIndices[j] = i

	// remove point
	c.liberties[c.numLiberties-1] = 0
	c.libertiesIndices[pt] = -1
	c.numLiberties--
}
