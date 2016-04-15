/*
Package board provides a library for placing stones on a Go game board.

It is inspired by 'Move Prediction in the Game of Go'. A thesis presented by Brett Alexander Harrison.
http://www.eecs.harvard.edu/econcs/pubs/Harrisonthesis.pdf
*/
package board

import (
	"errors"
)

type state int

const (
	black state = iota
	white
	empty
	wall
)

/*
A Board contains data of a Go board.

	7 by 7 board example.

	# # # # # # # #         00 01 02 03 04 05 06 07
	# . . . . . . .         08 09 10 11 12 13 14 15
	# . . . . . . .         16 17 18 19 20 21 22 23
	# . . . . . . .         24 25 26 27 28 29 30 31
	# . . . . . . .         32 33 34 35 36 37 38 39
	# . . . . . . .         40 41 42 43 44 45 46 47
	# . . . . . . .         48 49 50 51 52 53 54 55
	# . . . . . . .         56 57 58 59 60 61 62 63
	# # # # # # # #         64 65 66 67 68 69 70 71
	#                       72
*/
type Board struct {

	// boardSize = (size+2)*(size+1)+1
	size      int
	boardSize int

	// Max number of previous moves to store.
	maxHistory int

	// Array length is boardSize. chainReps - Zero if no chain.
	states    []state
	chains    []*chain
	chainReps []int

	// Current ko point if exists, 0 otherwise
	koPoint int

	// Number of stones captured
	blackDead int
	whiteDead int

	// Move history
	histories []*history
	depth     int
}

// NewBoard create a Board object.
func NewBoard(size int) Board {

	bh := Board{
		size:       size,
		maxHistory: 600,
	}
	bh.init()

	return bh
}

func (bd *Board) init() {

	bd.boardSize = (bd.size+2)*(bd.size+1) + 1

	// Index zero is not used.
	bd.histories = make([]*history, bd.maxHistory+1)

	bd.states = make([]state, bd.boardSize)

	bd.chains = make([]*chain, bd.boardSize)

	bd.chainReps = make([]int, bd.boardSize)

	bd.initStates()
}

func (bd *Board) initStates() {

	for i := 0; i <= bd.size+2; i++ {

		lead := i * (bd.size + 1)

		if i == 0 || i == bd.size+1 {

			for j := lead; j < lead+(bd.size+1); j++ {
				bd.states[j] = wall
			}

		} else if i == bd.size+2 {

			bd.states[lead] = wall

		} else {

			bd.states[lead] = wall

			for j := lead + 1; j < lead+(bd.size+1); j++ {
				bd.states[j] = empty
			}
		}
	}
}

// String is the text representation of current board state.
func (bd *Board) String() string {

	var r string

	for i, s := range bd.states {

		var c string

		switch s {
		case empty:
			c = "."
		case wall:
			c = "#"
		case black:
			c = "X"
		case white:
			c = "O"
		default:
			c = "?"
		}

		r += c

		if (i+1)%(bd.size+1) == 0 && i != 0 {
			r += "\n"
		}
	}

	return r
}

// DoBlack puts a black stone on a point.
func (bd *Board) DoBlack(pt int) error {

	return bd.do(pt, black)
}

// DoWhite puts a white stone on a point.
func (bd *Board) DoWhite(pt int) error {

	return bd.do(pt, white)
}

func (bd *Board) do(pt int, clr state) error {

	err := bd.isLegal(pt, clr)
	if err != nil {
		return err
	}

	h := newHistory(clr, pt, bd.koPoint)

	c := newChain(bd.size)
	c.addPoint(pt)

	// Initalize captured
	cp := newChain(bd.size)

	nb := bd.neighbors(pt)

	for i := 0; i < 4; i++ {

		n := nb[i]

		if bd.states[n] == empty {

			c.addLiberty(n)

		} else if bd.states[n] == clr && c.hasPoint(n) == false {

			c = *bd.joinChains(&c, bd.chains[n])

			bd.updateLibertiesAndChainReps(&c, clr)

		} else if bd.states[n] == bd.oppositePlayer(clr) {

			nc := bd.chains[n]

			if nc.numLiberties == 1 {

				bd.removeFromBoard(nc)

				bd.updatePrisoners(nc, clr)

				//Push
				for j := 0; j < nc.numPoints; j++ {

					ncp := nc.points[j]

					cp.addPoint(ncp)
				}

				bd.updateNeighboringChainsLiberties(nc)

				h.setCaptureDirections(i)
			}
		}
	}

	bd.updateLibertiesAndChainReps(&c, clr)

	bd.updateNeighboringChainsLiberties(&c)

	if cp.numPoints == 1 && c.numPoints == 1 {

		bd.koPoint = cp.points[0]

	} else {

		bd.koPoint = 0
	}

	bd.depth++

	bd.histories[bd.depth] = &h

	return nil
}

func (bd *Board) isLegal(pt int, clr state) error {

	if bd.depth >= bd.maxHistory {
		return errors.New("depth is larger than maxHistory")
	}

	if bd.isEmpty(pt) == false {
		return errors.New("point is not empty")
	}

	if bd.isKo(pt, clr) == true {
		return errors.New("point is Ko")
	}

	if bd.isSuicide(pt, clr) == true {
		return errors.New("point is suicide")
	}

	return nil
}

func (bd *Board) isEmpty(pt int) bool {

	return bd.states[pt] == empty
}

func (bd *Board) isKo(pt int, clr state) bool {

	result := false

	if pt == bd.koPoint {

		// This is for game ending winner fill in self ko.
		if bd.isAdjacentSelfChainWithTwoPlusLiberties(pt, clr) == false {
			result = true
		}
	}

	return result
}

func (bd *Board) isSuicide(pt int, clr state) bool {

	b1 := bd.isAdjacentEmpty(pt)

	b2 := bd.isAdjacentSelfChainWithTwoPlusLiberties(pt, clr)

	b3 := bd.isAdjacentEnemyChainWithOneLiberty(pt, clr)

	return !(b1 || b2 || b3)
}

func (bd *Board) isAdjacentSelfChainWithTwoPlusLiberties(pt int, clr state) bool {

	r := false

	nb := bd.neighbors(pt)

	for i := 0; i < 4; i++ {

		n := nb[i]

		if bd.states[n] == clr {

			if bd.chains[n].numLiberties >= 2 {

				r = true

				break
			}
		}
	}

	return r
}

func (bd *Board) isAdjacentEmpty(pt int) bool {

	nb := bd.neighbors(pt)

	return bd.states[nb[0]] == empty ||
		bd.states[nb[1]] == empty ||
		bd.states[nb[2]] == empty ||
		bd.states[nb[3]] == empty
}

func (bd *Board) isAdjacentEnemyChainWithOneLiberty(pt int, clr state) bool {

	r := false

	nb := bd.neighbors(pt)

	for i := 0; i < 4; i++ {

		n := nb[i]

		if bd.states[n] == bd.oppositePlayer(clr) {

			if bd.chains[n].numLiberties == 1 {

				r = true

				break
			}
		}
	}

	return r
}

func (bd *Board) oppositePlayer(clr state) state {

	r := clr

	if clr == black {
		r = white
	} else {
		r = black
	}

	return r
}

func (bd *Board) joinChains(c1 *chain, c2 *chain) *chain {

	// Add points and liberties of c2 to c1.
	for i := 0; i < c2.numPoints; i++ {
		c1.addPoint(c2.points[i])
	}

	return c1
}

func (bd *Board) updateLibertiesAndChainReps(c *chain, clr state) {

	for i := 0; i < c.numPoints; i++ {

		pt := c.points[i]

		// Update states, chains, chain_reps
		bd.states[pt] = clr

		bd.chains[pt] = c

		bd.chainReps[pt] = c.points[0]
	}

	for i := 0; i < c.numPoints; i++ {

		pt := c.points[i]

		nb := bd.neighbors(pt)

		for j := 0; j < 4; j++ {

			n := nb[j]

			if bd.states[n] == empty {
				c.addLiberty(n)
			}
		}
	}
}

func (bd *Board) removeFromBoard(c *chain) {

	for i := 0; i < c.numPoints; i++ {

		pt := c.points[i]

		bd.setEmpty(pt)
	}
}

func (bd *Board) updatePrisoners(nc *chain, clr state) {

	if clr == black {
		bd.blackDead += nc.numPoints
	} else if clr == white {
		bd.whiteDead += nc.numPoints
	}
}

func (bd *Board) updateNeighboringChainsLiberties(c *chain) {

	for i := 0; i < c.numPoints; i++ {

		pt := c.points[i]

		nb := bd.neighbors(pt)

		for j := 0; j < 4; j++ {

			n := nb[j]

			bd.updateLiberties(bd.chains[n])
		}
	}
}

func (bd *Board) setEmpty(pt int) {

	bd.states[pt] = empty
	bd.chains[pt] = nil
	bd.chainReps[pt] = 0
}

func (bd *Board) updateLiberties(c *chain) {

	if c == nil {
		return
	}

	for i := 0; i < c.numPoints; i++ {

		pt := c.points[i]

		nb := bd.neighbors(pt)

		for j := 0; j < 4; j++ {

			n := nb[j]

			if bd.states[n] == empty {

				c.addLiberty(n)

			} else {
				// This is needed for unknown Neighbors.
				c.removeLiberty(n)
			}
		}
	}
}

// Undo remove the last stone placed on the Go board.
func (bd *Board) Undo() error {

	if bd.depth == 0 {
		return errors.New("no history")
	}

	h := bd.histories[bd.depth]

	clr := h.color

	pt := h.point

	bd.setEmpty(pt)

	bd.koPoint = 0

	nb := bd.neighbors(pt)

	for i := 0; i < 4; i++ {

		n := nb[i]

		if bd.states[n] == bd.oppositePlayer(clr) {

			bd.chains[n].addLiberty(pt)

		} else if bd.states[n] == clr {

			c := bd.reconstructChain(n, clr, pt)

			bd.updateLibertiesAndChainReps(&c, clr)
		}

		if h.isCaptureDirections(i) == true {

			np := bd.oppositePlayer(clr)

			c := bd.reconstructChain(n, empty, pt)

			for j := 0; j < c.numPoints; j++ {
				bd.states[c.points[j]] = np
			}

			bd.updateLibertiesAndChainReps(&c, np)

			bd.updateNeighboringChainsLiberties(&c)

			// Update prisoners
			if clr == black {
				bd.blackDead -= c.numPoints
			} else if clr == white {
				bd.whiteDead -= c.numPoints
			}
		}
	}

	bd.koPoint = h.koPoint

	bd.depth--

	return nil
}

func (bd *Board) reconstructChain(pt int, clr state, original int) chain {

	c := newChain(bd.size)

	c.addPoint(pt)

	sps := bd.neighbors(pt)

	for len(sps) != 0 {

		len := len(sps)

		for i := len - 1; i >= 0; i-- {

			sp := sps[i]

			if bd.states[sp] == clr && c.hasPoint(sp) == false && sp != original {

				c.addPoint(sp)

				sps = append(sps, bd.neighbors(sp)...)
			}

			// remove sp
			sps = append(sps[:i], sps[i+1:]...)
		}
	}

	return c
}

// neighbors returns surrounding points with order north/east/south/west.
func (bd *Board) neighbors(pt int) []int {

	return []int{
		pt - (bd.size + 1),
		pt + 1,
		pt + (bd.size + 1),
		pt - 1}
}
