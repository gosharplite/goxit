package board

import (
	"testing"
)

func TestString(t *testing.T) {

	cases := map[string]struct {
		size     int
		expected string
	}{
		"3x3": {3, "####\n#...\n#...\n#...\n####\n"},
	}

	for k, tc := range cases {
		bh := New(tc.size)
		actual := bh.String()
		if actual != tc.expected {
			t.Errorf("%s: size %d,\n actual\n%s\n expected\n%s", k, tc.size, actual, tc.expected)
		}
	}
}

func TestDo(t *testing.T) {

	bh := New(3)

	bh.DoBlack(5)
	bh.DoWhite(6)
	bh.DoBlack(10)
	bh.DoWhite(9)

	actual := bh.String()
	expected := "####\n#.O.\n#OX.\n#...\n####\n"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}

	bh.Undo()

	actual = bh.String()
	expected = "####\n#XO.\n#.X.\n#...\n####\n"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}
}
