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
