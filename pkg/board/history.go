package board

type history struct {

	// Data to be recorded.
	color state
	point int

	// Ko point before move was played
	koPoint int

	// capture directions[d] = true if and only if
	// a capture occurred in the direction d from point
	captureDirections []bool
}

func newHistory(clr state, pt int, koPoint int) history {

	h := history{}

	h.color = clr

	h.point = pt

	h.koPoint = koPoint

	h.captureDirections = []bool{false, false, false, false}

	return h
}

func (h *history) setCaptureDirections(dir int) {

	h.captureDirections[dir] = true
}

func (h *history) isCaptureDirections(dir int) bool {

	return h.captureDirections[dir]
}
