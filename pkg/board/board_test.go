package board

import (
	"testing"
)

func TestString(t *testing.T) {

	cases := map[string]struct {
		size     int
		expected string
	}{
		"3x3": {3, "####\n#...\n#...\n#...\n####\n#"},
	}

	for k, tc := range cases {
		bh := NewBoard(tc.size)
		actual := bh.String()
		if actual != tc.expected {
			t.Errorf("%s: size %d,\n actual\n%s\n expected\n%s", k, tc.size, actual, tc.expected)
		}
	}
}

func TestMaxHistory(t *testing.T) {

	bh := Board{
		size:       3,
		maxHistory: 1,
	}
	bh.init()

	err := bh.DoBlack(5)
	if err != nil {
		t.Error(err.Error())
	}

	err = bh.DoWhite(6)
	if err == nil {
		t.Error("maxHistory exceeded but not detected")
	}
}

func TestIsEmpty(t *testing.T) {

	bh := NewBoard(3)

	err := bh.DoBlack(5)
	if err != nil {
		t.Error(err.Error())
	}

	err = bh.DoWhite(5)
	if err == nil {
		t.Error("point is not empty but not detected")
	}
}

func TestIsKo(t *testing.T) {

	bh := NewBoard(3)

	bh.DoBlack(5)
	bh.DoBlack(7)
	bh.DoBlack(10)

	bh.DoWhite(9)
	bh.DoWhite(6)

	err := bh.DoBlack(5)
	if err == nil {
		t.Error("point is Ko but not detected")
	}
}

func TestIsSuicide(t *testing.T) {

	bh := NewBoard(3)

	bh.DoBlack(6)
	bh.DoBlack(9)

	err := bh.DoWhite(5)
	if err == nil {
		t.Error("point is suicide but not detected")
	}
}

func TestDo(t *testing.T) {

	bh := NewBoard(3)

	bh.DoBlack(5)
	bh.DoWhite(6)
	bh.DoBlack(10)
	bh.DoWhite(9)

	actual := bh.String()
	expected := "####\n#.O.\n#OX.\n#...\n####\n#"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}

	bh.Undo()

	actual = bh.String()
	expected = "####\n#XO.\n#.X.\n#...\n####\n#"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}
}

func TestCapture1(t *testing.T) {

	bh := NewBoard(3)

	bh.DoBlack(5)
	bh.DoBlack(6)
	bh.DoBlack(9)
	bh.DoBlack(10)

	bh.DoWhite(7)
	bh.DoWhite(11)
	bh.DoWhite(13)
	bh.DoWhite(15)

	bh.DoWhite(14)

	actual := bh.String()
	expected := "####\n#..O\n#..O\n#OOO\n####\n#"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}

	bh.Undo()

	actual = bh.String()
	expected = "####\n#XXO\n#XXO\n#O.O\n####\n#"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}
}

func TestCapture2(t *testing.T) {

	bh := NewBoard(7)

	bh.DoBlack(26)
	bh.DoBlack(33)
	bh.DoBlack(35)
	bh.DoBlack(36)
	bh.DoBlack(42)

	bh.DoWhite(18)
	bh.DoWhite(25)
	bh.DoWhite(27)
	bh.DoWhite(28)
	bh.DoWhite(37)
	bh.DoWhite(41)
	bh.DoWhite(43)
	bh.DoWhite(44)
	bh.DoWhite(50)

	bh.DoWhite(34)

	actual := bh.String()
	expected := "########\n#.......\n#.O.....\n#O.OO...\n#.O..O..\n#O.OO...\n#.O.....\n#.......\n########\n#"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}

	bh.Undo()

	actual = bh.String()
	expected = "########\n#.......\n#.O.....\n#OXOO...\n#X.XXO..\n#OXOO...\n#.O.....\n#.......\n########\n#"
	if actual != expected {
		t.Errorf("\n actual\n%s\n expected\n%s", actual, expected)
	}
}

var result Board

func BenchmarkCapture(b *testing.B) {

	var bh Board

	for n := 0; n < b.N; n++ {
		bh = NewBoard(3)

		bh.DoBlack(5)
		bh.DoBlack(6)
		bh.DoBlack(9)
		bh.DoBlack(10)

		bh.DoWhite(7)
		bh.DoWhite(11)
		bh.DoWhite(13)
		bh.DoWhite(15)

		bh.DoWhite(14)
		bh.Undo()
	}

	result = bh
}
