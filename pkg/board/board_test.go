package board

import (
	"testing"
)

func TestLogBoard(t *testing.T) {

	bh, err := NewBoard(3)
	if err != nil {
		t.Error("Test failed")
	}

	expected := "####\n#...\n#...\n#...\n####\n"
	actual := bh.LogBoard()
	if actual != expected {
		t.Error("Test failed")
	}
}
