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

type Board struct {

	// Size of Go board.
	// boardSize = (size+2)*(size+1)+1
	size      int
	boardSize int

	// Max number of previous moves to store.
	maxHistory int

	// Arrays for storing states, chains, and chain representatives.
	// Array length is boardSize.
	// chainReps - Zero if no chain.
	states    []state
	chains    []*Chain
	chainReps []int

	// Current ko point if exists, 0 otherwise
	koPoint int

	// Number of stones captured from each player
	blackDead int
	whiteDead int

	// Move history list
	histories []*History
	depth     int
}

func New(size int) *Board {

	bh := Board{}
	bh.init(size)
	return &bh
}

func (bd *Board) init(size int) {

	bd.size = size
	bd.boardSize = (bd.size+2)*(bd.size+1) + 1

	bd.maxHistory = 600

	// Index zero is not used.
	bd.histories = make([]*History, bd.maxHistory+1)

	bd.states = make([]state, bd.boardSize)

	bd.chains = make([]*Chain, bd.boardSize)

	bd.chainReps = make([]int, bd.boardSize)

	bd.initStates()
}

func (bd *Board) initStates() {

	for i := 0; i <= bd.size+2; i++ {

		leadPosition := i * (bd.size + 1)

		if i == 0 || i == bd.size+1 {

			for j := leadPosition; j < leadPosition+(bd.size+1); j++ {
				bd.states[j] = wall
			}

		} else if i == bd.size+2 {

			bd.states[leadPosition] = wall

		} else {

			bd.states[leadPosition] = wall

			for j := leadPosition + 1; j < leadPosition+(bd.size+1); j++ {
				bd.states[j] = empty
			}
		}
	}

	//              7 by 7 example.
	//
	//              # # # # # # # #         00 01 02 03 04 05 06 07
	//              # . . . . . . .         08 09 10 11 12 13 14 15
	//              # . . . . . . .         16 17 18 19 20 21 22 23
	//              # . . . . . . .         24 25 26 27 28 29 30 31
	//              # . . . . . . .         32 33 34 35 36 37 38 39
	//              # . . . . . . .         40 41 42 43 44 45 46 47
	//              # . . . . . . .         48 49 50 51 52 53 54 55
	//              # . . . . . . .         56 57 58 59 60 61 62 63
	//              # # # # # # # #         64 65 66 67 68 69 70 71
	//              #                       72
}

func (bd *Board) String() string {

	var line, result string

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

		if i%(bd.size+1) == 0 && i != 0 {
			result += line + "\n"
			line = c
		} else {

			line += c
		}
	}

	return result
}

func (bd *Board) DoBlack(pt int) error {

	return bd.do(pt, black)
}

func (bd *Board) DoWhite(pt int) error {

	return bd.do(pt, white)
}

func (bd *Board) do(pt int, clr state) error {

	err := bd.isLegal(pt, clr)
	if err != nil {
		return err
	}

	// Initialize history
	history := History{}
	history.Init(clr, pt, bd.koPoint)

	// Initialize chain
	chain := Chain{}
	chain.Init(bd.size)
	chain.addPoint(pt)

	captured := Chain{}
	captured.Init(bd.size)

	// Same order as Direction in MoveHistory.
	neighbors := []int{bd.north(pt), bd.east(pt), bd.south(pt), bd.west(pt)}

	for i := 0; i < 4; i++ {

		n := neighbors[i]

		if bd.states[n] == empty {

			chain.addLiberty(n)

		} else if bd.states[n] == clr && chain.hasPoint(n) == false {

			chain = *bd.joinChains(&chain, bd.chains[n])

			bd.updateLibertiesAndChainReps(&chain, clr)

		} else if bd.states[n] == bd.oppositePlayer(clr) {

			nc := bd.chains[n]

			if nc.numLiberties == 1 {

				bd.removeFromBoard(nc)

				bd.updatePrisoners(nc, clr)

				//Push
				for j := 0; j < nc.numPoints; j++ {

					ncp := nc.points[j]

					captured.addPoint(ncp)
				}

				bd.updateNeighboringChainsLiberties(nc)

				history.setCaptureDirections(i)
			}
		}
	}

	bd.updateLibertiesAndChainReps(&chain, clr)

	bd.updateNeighboringChainsLiberties(&chain)

	if captured.numPoints == 1 && chain.numPoints == 1 {

		bd.koPoint = captured.points[0]

	} else {

		bd.koPoint = 0
	}

	bd.depth++

	bd.histories[bd.depth] = &history

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

	bisAdjacentEmpty := bd.isAdjacentEmpty(pt)

	bisAdjacentSelfChainWithTwoPlusLiberties := bd.isAdjacentSelfChainWithTwoPlusLiberties(pt, clr)

	bisAdjacentEnemyChainWithOneLiberty := bd.isAdjacentEnemyChainWithOneLiberty(pt, clr)

	return !(bisAdjacentEmpty ||
		bisAdjacentSelfChainWithTwoPlusLiberties ||
		bisAdjacentEnemyChainWithOneLiberty)
}

func (bd *Board) isAdjacentSelfChainWithTwoPlusLiberties(pt int, clr state) bool {

	result := false

	neighbors := []int{bd.north(pt), bd.south(pt), bd.east(pt), bd.west(pt)}

	for i := 0; i < 4; i++ {

		n := neighbors[i]

		if bd.states[n] == clr {

			if bd.chains[n].numLiberties >= 2 {

				result = true

				break
			}
		}
	}

	return result
}

func (bd *Board) isAdjacentEmpty(pt int) bool {

	return bd.states[bd.north(pt)] == empty ||
		bd.states[bd.south(pt)] == empty ||
		bd.states[bd.east(pt)] == empty ||
		bd.states[bd.west(pt)] == empty
}

func (bd *Board) isAdjacentEnemyChainWithOneLiberty(pt int, clr state) bool {

	result := false

	neighbors := []int{bd.north(pt), bd.south(pt), bd.east(pt), bd.west(pt)}

	for i := 0; i < 4; i++ {

		n := neighbors[i]

		if bd.states[n] == bd.oppositePlayer(clr) {

			if bd.chains[n].numLiberties == 1 {

				result = true

				break
			}
		}
	}

	return result
}

func (bd *Board) north(pt int) int {

	return pt - (bd.size + 1)
}

func (bd *Board) south(pt int) int {

	return pt + (bd.size + 1)
}

func (bd *Board) east(pt int) int {

	return pt + 1
}

func (bd *Board) west(pt int) int {

	return pt - 1
}

func (bd *Board) oppositePlayer(clr state) state {

	result := clr

	if clr == black {
		result = white
	} else {
		result = black
	}

	return result
}

func (bd *Board) joinChains(c1 *Chain, c2 *Chain) *Chain {

	// Add points and liberties of c2 to c1.
	for i := 0; i < c2.numPoints; i++ {
		c1.addPoint(c2.points[i])
	}

	return c1
}

func (bd *Board) updateLibertiesAndChainReps(c *Chain, clr state) {

	for i := 0; i < c.numPoints; i++ {

		point := c.points[i]

		// Update states, chains, chain_reps
		bd.states[point] = clr

		bd.chains[point] = c

		bd.chainReps[point] = c.points[0]
	}

	for i := 0; i < c.numPoints; i++ {

		pt := c.points[i]

		// Update liberties.
		neighbors := []int{bd.north(pt), bd.south(pt), bd.east(pt), bd.west(pt)}

		for j := 0; j < 4; j++ {

			n := neighbors[j]

			if bd.states[n] == empty {
				c.addLiberty(n)
			}
		}
	}
}

func (bd *Board) removeFromBoard(c *Chain) {

	for i := 0; i < c.numPoints; i++ {

		pt := c.points[i]

		bd.setEmpty(pt)
	}
}

func (bd *Board) updatePrisoners(nc *Chain, clr state) {

	if clr == black {
		bd.blackDead += nc.numPoints
	} else if clr == white {
		bd.whiteDead += nc.numPoints
	}
}

func (bd *Board) updateNeighboringChainsLiberties(c *Chain) {

	for i := 0; i < c.numPoints; i++ {

		pt := c.points[i]

		// Update liberties.
		neighbors := []int{bd.north(pt), bd.south(pt), bd.east(pt), bd.west(pt)}

		for j := 0; j < 4; j++ {

			n := neighbors[j]

			bd.updateLiberties(bd.chains[n])
		}
	}
}

func (bd *Board) setEmpty(pt int) {

	bd.states[pt] = empty
	bd.chains[pt] = nil
	bd.chainReps[pt] = 0
}

func (bd *Board) updateLiberties(c *Chain) {

	if c == nil {
		return
	}

	for i := 0; i < c.numPoints; i++ {

		pt := c.points[i]

		// Update liberties.
		neighbors := []int{bd.north(pt), bd.south(pt), bd.east(pt), bd.west(pt)}

		for j := 0; j < 4; j++ {

			n := neighbors[j]

			if bd.states[n] == empty {

				c.addLiberty(n)

			} else {

				c.removeLiberty(n) // This is needed for unknown Neighbors.
			}
		}
	}
}

func (bd *Board) Undo() error {

	if bd.depth == 0 {
		return errors.New("no history")
	}

	h := bd.histories[bd.depth]

	clr := h.color

	pt := h.point

	bd.setEmpty(pt)

	bd.koPoint = 0

	// Same order as Direction in History.
	neighbors := []int{bd.north(pt), bd.east(pt), bd.south(pt), bd.west(pt)}

	for i := 0; i < 4; i++ {

		n := neighbors[i]

		if bd.states[n] == bd.oppositePlayer(clr) {

			bd.chains[n].addLiberty(pt)

		} else if bd.states[n] == clr {

			chain := bd.reconstructChain(n, clr, pt)

			bd.updateLibertiesAndChainReps(chain, clr)
		}

		if h.isCaptureDirections(i) == true {

			np := bd.oppositePlayer(clr)

			chain := bd.reconstructChain(n, empty, pt)

			for j := 0; j < chain.numPoints; j++ {
				bd.states[chain.points[j]] = np
			}

			bd.updateLibertiesAndChainReps(chain, np)

			bd.updateNeighboringChainsLiberties(chain)

			// Update prisoners
			if clr == black {
				bd.blackDead -= chain.numPoints
			} else if clr == white {
				bd.whiteDead -= chain.numPoints
			}
		}
	}

	bd.koPoint = h.koPoint

	bd.depth--

	return nil
}

func (bd *Board) reconstructChain(pt int, clr state, original int) *Chain {

	chain := Chain{}
	chain.Init(bd.size)
	chain.addPoint(pt)

	searchPoints := bd.getneighbors(pt)

	for len(searchPoints) != 0 {

		len := len(searchPoints)

		for i := len - 1; i >= 0; i-- {

			sp := searchPoints[i]

			if bd.states[sp] == clr && chain.hasPoint(sp) == false && sp != original {

				chain.addPoint(sp)

				searchPoints = append(searchPoints, bd.getneighbors(sp)...)
			}

			// remove sp
			front := searchPoints[:i]
			back := searchPoints[i+1:]

			searchPoints = append(front, back...)
		}
	}

	return &chain
}

func (bd *Board) getneighbors(pt int) []int {

	result := make([]int, 0)

	result = append(result,
		bd.north(pt),
		bd.east(pt),
		bd.south(pt),
		bd.west(pt))

	return result
}
