package board

import (
	"testing"
)

func TestLogBoard(t *testing.T) {

	bh := NewBoard(3)

	expected := "####\n#...\n#...\n#...\n####\n"
	actual := bh.LogBoard()
	if actual != expected {
		t.Error("Test failed")
	}
}
