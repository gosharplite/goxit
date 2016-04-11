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

func TestProcessMoveUndoMove(t *testing.T) {

	bh := NewBoard(3)

	bh.ProcessMove(5, State_BLACK)
	bh.ProcessMove(6, State_WHITE)
	bh.ProcessMove(10, State_BLACK)
	bh.ProcessMove(9, State_WHITE)

	actual := bh.LogBoard()
	expected := "####\n#.O.\n#OX.\n#...\n####\n"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}

	bh.UndoMove()

	actual = bh.LogBoard()
	expected = "####\n#XO.\n#.X.\n#...\n####\n"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}
}
