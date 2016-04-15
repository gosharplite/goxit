package hash

import (
	"testing"
)

func TestCanonical(t *testing.T) {

	p := NewPattern(5)

	p.SetBlack(3, 0)
	p.SetBlack(2, 1)
	p.SetBlack(2, 2)
	p.SetWhite(4, 2)
	p.SetWhite(3, 3)

	p.canonical()
	h1 := p.GetHash()

	p = Pattern{}
	p.init(5)

	p.SetWhite(1, 2)
	p.SetWhite(2, 2)
	p.SetWhite(0, 3)
	p.SetBlack(3, 3)
	p.SetBlack(2, 4)

	p.canonical()

	h2 := p.GetHash()

	if h1 != h2 {
		t.Errorf("%v != %v", h1, h2)
	}

	p.string(p.black, p.white)
}
