package board

import ()

type BoardState int

const (
	// State of a board intersection
	State_BLACK BoardState = 1
	State_WHITE BoardState = 2
	State_EMPTY BoardState = 3
	State_WALL  BoardState = 4
)

type BoardHarvard struct {

	// Size parameters for a Go board
	SIZE       int
	BOARD_SIZE int //(SIZE+2)*(SIZE+1)+1;

	// Max number of previous moves to store
	MAX_HISTORY int

	// Arrays for storing states, chains, and chain representatives
	states     []BoardState //BOARD_SIZE
	chains     []*Chain     //BOARD_SIZE
	chain_reps []int        //BOARD_SIZE. Zero if no chain.

	// Current ko point if exists, 0 otherwise
	ko_point int

	// Number of stones captured from each player
	black_prisoners int
	white_prisoners int

	// Move history list
	move_history_list []*MoveHistory
	Depth             int
}

func NewBoard(size int) *BoardHarvard {
	bh := BoardHarvard{}
	bh.Initialize(size)
	return &bh
}

func (board *BoardHarvard) Initialize(size int) {

	board.SIZE = size
	board.BOARD_SIZE = (board.SIZE+2)*(board.SIZE+1) + 1

	board.MAX_HISTORY = 600

	board.move_history_list = make([]*MoveHistory, board.MAX_HISTORY)

	board.states = make([]BoardState, board.BOARD_SIZE)

	board.chains = make([]*Chain, board.BOARD_SIZE)

	board.chain_reps = make([]int, board.BOARD_SIZE)

	board.initializeStates()
}

func (board *BoardHarvard) initializeStates() {

	for i := 0; i <= board.SIZE+2; i++ {

		leadPosition := i * (board.SIZE + 1)

		if i == 0 || i == board.SIZE+1 {

			for j := leadPosition; j < leadPosition+(board.SIZE+1); j++ {
				board.states[j] = State_WALL
			}

		} else if i == board.SIZE+2 {

			board.states[leadPosition] = State_WALL

		} else {

			board.states[leadPosition] = State_WALL

			for j := leadPosition + 1; j < leadPosition+(board.SIZE+1); j++ {
				board.states[j] = State_EMPTY
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

func (board *BoardHarvard) LogBoard() string {

	var line, result string

	for i, s := range board.states {

		var c string

		switch s {
		case State_EMPTY:
			c = "."
		case State_WALL:
			c = "#"
		case State_BLACK:
			c = "X"
		case State_WHITE:
			c = "O"
		default:
			c = "?"
		}

		if i%(board.SIZE+1) == 0 && i != 0 {
			result += line + "\n"
			line = c
		} else {

			line += c
		}
	}

	return result
}

func (board *BoardHarvard) ProcessMove(point int, player BoardState) bool {

	isLegal := true

	if board.isLegalMove(point, player) == false {
		return false
	}

	// Initialize MoveHistory
	moveHistory := MoveHistory{}
	moveHistory.Initialize(player, point, board.ko_point)

	// Initialize new chain
	chain := Chain{}
	chain.Initialize(board.SIZE)
	chain.addPoint(point)

	captured := Chain{}
	captured.Initialize(board.SIZE)

	// Same order as Direction in MoveHistory.
	neighbors := []int{board.north(point), board.east(point), board.south(point), board.west(point)}

	for i := 0; i < 4; i++ {

		n := neighbors[i]

		if board.states[n] == State_EMPTY {

			chain.addLiberty(n)

		} else if board.states[n] == player && chain.hasPoint(n) == false {

			chain = *board.joinChains(&chain, board.chains[n])

			board.updateLibertiesAndChainReps(&chain, player)

		} else if board.states[n] == board.oppositePlayer(player) {

			nc := board.chains[n]

			if nc.num_liberties == 1 {

				board.removeFromBoard(nc)

				board.updatePrisoners(nc, player)

				//Push
				for j := 0; j < nc.num_points; j++ {

					ncp := nc.points[j]

					captured.addPoint(ncp)
				}
				board.updateNeighboringChainsLiberties(nc)

				moveHistory.setCapture_directions(i)
			}
		}
	}

	board.updateLibertiesAndChainReps(&chain, player)
	board.updateNeighboringChainsLiberties(&chain)

	if captured.num_points == 1 && chain.num_points == 1 { // Paper is wrong about Ko.

		board.ko_point = captured.points[0]

	} else {

		board.ko_point = 0
	}

	board.Depth++

	board.move_history_list[board.Depth] = &moveHistory

	return isLegal
}

func (board *BoardHarvard) isLegalMove(point int, player BoardState) bool {

	return board.isEmpty(point) == true &&
		board.isKo(point, player) == false &&
		board.isNotSuicide(point, player) == true

	return false
}

func (board *BoardHarvard) isEmpty(point int) bool {
	return board.states[point] == State_EMPTY
}

func (board *BoardHarvard) isKo(point int, player BoardState) bool {

	result := false

	if point == board.ko_point {

		// This is for game ending winner fill in self ko.
		if board.isAdjacentSelfChainWithTwoPlusLiberties(point, player) == false {
			result = true
		}
	}

	return result
}

func (board *BoardHarvard) isNotSuicide(point int, player BoardState) bool {

	bisAdjacentEmpty := board.isAdjacentEmpty(point)

	bisAdjacentSelfChainWithTwoPlusLiberties := board.isAdjacentSelfChainWithTwoPlusLiberties(point, player)

	bisAdjacentEnemyChainWithOneLiberty := board.isAdjacentEnemyChainWithOneLiberty(point, player)

	return bisAdjacentEmpty ||
		bisAdjacentSelfChainWithTwoPlusLiberties ||
		bisAdjacentEnemyChainWithOneLiberty

	return false
}

func (board *BoardHarvard) isAdjacentSelfChainWithTwoPlusLiberties(point int, player BoardState) bool {

	result := false

	neighbors := []int{board.north(point), board.south(point), board.east(point), board.west(point)}

	for i := 0; i < 4; i++ {

		n := neighbors[i]

		if board.states[n] == player {

			if board.chains[n].num_liberties >= 2 {

				result = true

				break
			}
		}
	}

	return result
}

func (board *BoardHarvard) isAdjacentEmpty(point int) bool {

	return board.states[board.north(point)] == State_EMPTY ||
		board.states[board.south(point)] == State_EMPTY ||
		board.states[board.east(point)] == State_EMPTY ||
		board.states[board.west(point)] == State_EMPTY
}

func (board *BoardHarvard) isAdjacentEnemyChainWithOneLiberty(point int, player BoardState) bool {

	result := false

	neighbors := []int{board.north(point), board.south(point), board.east(point), board.west(point)}

	for i := 0; i < 4; i++ {

		n := neighbors[i]

		if board.states[n] == board.oppositePlayer(player) {

			if board.chains[n].num_liberties == 1 {

				result = true

				break
			}
		}
	}

	return result
}

func (board *BoardHarvard) north(point int) int {
	return point - (board.SIZE + 1)
}

func (board *BoardHarvard) south(point int) int {
	return point + (board.SIZE + 1)
}

func (board *BoardHarvard) east(point int) int {
	return point + 1
}

func (board *BoardHarvard) west(point int) int {
	return point - 1
}

func (board *BoardHarvard) oppositePlayer(player BoardState) BoardState {

	result := player

	if player == State_BLACK {
		result = State_WHITE
	} else {
		result = State_BLACK
	}

	return result
}

func (board *BoardHarvard) joinChains(c1 *Chain, c2 *Chain) *Chain {

	// Add points and liberties of c2 to c1.
	for i := 0; i < c2.num_points; i++ {
		c1.addPoint(c2.points[i])
	}

	return c1
}

func (board *BoardHarvard) updateLibertiesAndChainReps(chain *Chain, player BoardState) {

	for i := 0; i < chain.num_points; i++ {

		point := chain.points[i]

		// Update states, chains, chain_reps
		board.states[point] = player

		board.chains[point] = chain

		board.chain_reps[point] = chain.points[0]
	}

	for i := 0; i < chain.num_points; i++ {

		point := chain.points[i]

		// Update liberties.
		neighbors := []int{board.north(point), board.south(point), board.east(point), board.west(point)}

		for j := 0; j < 4; j++ {

			n := neighbors[j]

			if board.states[n] == State_EMPTY {
				chain.addLiberty(n)
			}
		}
	}
}

func (board *BoardHarvard) removeFromBoard(chain *Chain) {

	//slog.Debug("num_points: %v", chain.num_points)

	for i := 0; i < chain.num_points; i++ {

		//slog.Debug("i: %v", i)

		point := chain.points[i]

		//slog.Debug("point: %v, %v", i, point)

		board.setEmpty(point)
	}
}

func (board *BoardHarvard) updatePrisoners(nc *Chain, player BoardState) {

	if player == State_BLACK {
		board.black_prisoners += nc.num_points
	} else if player == State_WHITE {
		board.white_prisoners += nc.num_points
	}
}

func (board *BoardHarvard) updateNeighboringChainsLiberties(chain *Chain) {

	for i := 0; i < chain.num_points; i++ {

		point := chain.points[i]

		// Update liberties.
		neighbors := []int{board.north(point), board.south(point), board.east(point), board.west(point)}

		for j := 0; j < 4; j++ {

			n := neighbors[j]

			board.updateLiberties(board.chains[n])
		}
	}
}

func (board *BoardHarvard) setEmpty(point int) {

	board.states[point] = State_EMPTY
	board.chains[point] = nil
	board.chain_reps[point] = 0
}

func (board *BoardHarvard) updateLiberties(chain *Chain) {

	if chain == nil {
		return
	}

	for i := 0; i < chain.num_points; i++ {

		point := chain.points[i]

		// Update liberties.
		neighbors := []int{board.north(point), board.south(point), board.east(point), board.west(point)}

		for j := 0; j < 4; j++ {

			n := neighbors[j]

			if board.states[n] == State_EMPTY {

				chain.addLiberty(n)

			} else {

				chain.removeLiberty(n) // This is needed for unknown Neighbors.
			}
		}
	}
}
