package board

import "testing"

func TestLogBoard(t *testing.T) {

	bh := BoardHarvard{}

	bh.Initialize(9)

	bh.LogBoard()

	//	expected := "Hello Go!"
	//	actual := hello()
	//	if actual != expected {
	//		t.Error("Test failed")
	//	}
}
