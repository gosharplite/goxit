package board

import (
	"testing"
)

func TestLogBoard(t *testing.T) {

	cases := map[string]struct {
		size     int
		expected string
	}{
		"3x3": {3, "####\n#...\n#...\n#...\n####\n"},
	}

	for k, tc := range cases {
		bh := NewBoard(tc.size)
		actual := bh.LogBoard()
		if actual != tc.expected {
			t.Errorf("%s: size %d,\n actual\n%s\n expected\n%s", k, tc.size, actual, tc.expected)
		}
	}
}

func TestProcessMove(t *testing.T) {

	type stone struct {
		point int
		state BoardState
	}

	cases := map[string]struct {
		size     int
		points   stone
		expected stone
	}{
		"simple point": {3, stone{5, State_BLACK}, stone{5, State_BLACK}},
	}
}
