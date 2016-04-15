package hash

import (
	"testing"
)

func TestCanonical(t *testing.T) {

	p := Pattern{}
	p.Initialize(5)

	p.SetBlack(3, 0)
	p.SetBlack(2, 1)
	p.SetBlack(2, 2)
	p.SetWhite(4, 2)
	p.SetWhite(3, 3)

	p.Canonical()
	h1 := p.GetHashMax()

	p = Pattern{}
	p.Initialize(5)

	p.SetWhite(1, 2)
	p.SetWhite(2, 2)
	p.SetWhite(0, 3)
	p.SetBlack(3, 3)
	p.SetBlack(2, 4)

	p.Canonical()
	h2 := p.GetHashMax()

	if h1 != h2 {
		t.Errorf("%v != %v", h1, h2)
	}
}
