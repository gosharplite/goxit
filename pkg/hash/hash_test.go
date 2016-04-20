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

	p := NewPattern(5)

	p.SetBlack(3, 0)
	p.SetBlack(2, 1)
	p.SetBlack(2, 2)
	p.SetWhite(4, 2)
	p.SetWhite(3, 3)

	h1 := p.GetHash()

	p = Pattern{}
	p.init(5)

	p.SetWhite(1, 2)
	p.SetWhite(2, 2)
	p.SetWhite(0, 3)
	p.SetBlack(3, 3)
	p.SetBlack(2, 4)

	h2 := p.GetHash()

	if h1 != h2 {
		t.Errorf("%v != %v", h1, h2)
	}

	if h1 != 11305780 {
		t.Errorf("%v != 11305780", h1)
	}

	//fmt.Printf("%v\n", h1)
	//fmt.Printf("%v", p.string(p.black, p.white))
}

func BenchmarkPerformance(b *testing.B) {

	p := NewPattern(19)

	for n := 0; n < b.N; n++ {
		p.performance()
	}
}

var result int64

func BenchmarkGetHash(b *testing.B) {

	var h int64

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
