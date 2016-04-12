package board

import ()

type History struct {

	// Data to be recorded.
	color state
	point int

	// Ko point before move was played
	koPoint int

	// capture directions[d] = true if and only if
	// a capture occurred in the direction d from point
	captureDirections []bool
}

func (h *History) Init(clr state, pt int, koPoint int) {

	h.color = clr

	h.point = pt

	h.koPoint = koPoint

	h.captureDirections = []bool{false, false, false, false}
}

func (h *History) setCaptureDirections(dir int) {

	h.captureDirections[dir] = true
}

func (h *History) isCaptureDirections(dir int) bool {

	return h.captureDirections[dir]
}
