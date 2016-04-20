package hash

import (
	"testing"
)

func TestTime(t *testing.T) {

	p := NewPattern(19)

	p.SetBlack(3, 0)
	p.SetBlack(2, 1)
	p.SetBlack(2, 2)
	p.SetWhite(4, 2)
	p.SetWhite(3, 3)

	p.GetHash()
}

func TestCanonical(t *testing.T) {

	p := NewPattern(19)

	p.SetBlack(3, 0)
	p.SetBlack(2, 1)
	p.SetBlack(2, 2)
	p.SetWhite(4, 2)
	p.SetWhite(3, 3)

	h1 := p.GetHash()

	p = Pattern{}
	p.init(19)

	p.SetWhite(18, 15)
	p.SetWhite(17, 16)
	p.SetWhite(16, 16)
	p.SetBlack(16, 14)
	p.SetBlack(15, 15)

	h2 := p.GetHash()

	if h1 != h2 {
		t.Errorf("%v != %v", h1, h2)
	}

	if h1 != 79471947460759914 {
		t.Errorf("%v != 79471947460759914", h1)
	}
}

var result uint64

func BenchmarkGetHash(b *testing.B) {

	var h uint64

	p := NewPattern(19)

	p.SetBlack(3, 0)
	p.SetBlack(2, 1)
	p.SetBlack(2, 2)
	p.SetWhite(4, 2)
	p.SetWhite(3, 3)

	for n := 0; n < b.N; n++ {
		h = p.GetHash()
	}

	result = h
}
