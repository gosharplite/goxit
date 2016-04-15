/*
Package hash provides a library for calculating a canonical hash value of a Go game pattern.

Canonical hash value means same value will be found if the pattern is mirrored, rotated or color changed. Below patterns will have same hash value.

	...X.	...O.	.O...	.....
	..X..	..O..	..O..	.....
	..X.O	..O.X	X.O..	.OO..
	...O.	...X.	.X...	O..X.
	.....	.....	.....	..X..
*/
package hash

import (
	"fmt"
	"github.com/willf/bitset"
)

// A Pattern contains data of a Go game situation.
type Pattern struct {
	size int

	black *bitset.BitSet
	white *bitset.BitSet

	blacks []*bitset.BitSet
	whites []*bitset.BitSet

	mirror []int
	rot90  []int
	rot180 []int
	rot270 []int

	prime  int64
	celing int64
}

func NewPattern(size int) Pattern {

	p := Pattern{}
	p.init(size)

	return p
}

func (p *Pattern) init(size int) {

	p.size = size

	p.black = bitset.New(uint(size * size))

	p.white = bitset.New(uint(size * size))

	p.mirror = []int{1, 0, 0, -1}
	p.rot90 = []int{0, -1, 1, 0}
	p.rot180 = []int{-1, 0, 0, -1}
	p.rot270 = []int{0, 1, -1, 0}

	p.prime = 13
	p.celing = (9223372036854775807 - 4294967295) / p.prime
}

func (p *Pattern) SetBlack(x, y int) {

	q := uint(y*p.size + x)

	p.black.Set(q)
}

func (p *Pattern) SetWhite(x, y int) {

	q := uint(y*p.size + x)

	p.white.Set(q)
}

func (p *Pattern) Canonical() {

	l := uint(p.size * p.size)

	p.blacks = p.initBitSet(p.black, l)
	p.whites = p.initBitSet(p.white, l)

	for i := uint(0); i < l; i++ {

		if p.black.Test(i) {
			p.translate(i, p.blacks)
		}

		if p.white.Test(i) {
			p.translate(i, p.whites)
		}
	}

	// Reverse color
	for i := 8; i < 16; i++ {

		p.blacks[i] = p.whites[i-8]

		p.whites[i] = p.blacks[i-8]
	}

	// Find canonical
	bc := p.black
	wc := p.white

	bwc := p.combineBitSet(l, bc, wc)

	for j := 1; j < 16; j++ {

		a := p.combineBitSet(l, p.blacks[j], p.whites[j])

		b := bitset.New(l * 2)
		b = b.Union(bwc)

		// Find different bits.
		b = b.SymmetricDifference(a)

		// Check first true bit.
		n := p.nextSetBit(0, b)
		if n >= 0 {

			if bwc.Test(uint(n)) == false {

				bc = p.blacks[j]
				wc = p.whites[j]

				bwc = a
			}
		}
	}

	p.black = bc
	p.white = wc
}

func (p *Pattern) initBitSet(bitSet *bitset.BitSet, length uint) []*bitset.BitSet {

	result := []*bitset.BitSet{
		bitSet, // 0
		bitset.New(length),
		bitset.New(length),
		bitset.New(length),
		bitset.New(length),
		bitset.New(length),
		bitset.New(length),
		bitset.New(length), // 7
		bitset.New(length),
		bitset.New(length),
		bitset.New(length),
		bitset.New(length),
		bitset.New(length),
		bitset.New(length),
		bitset.New(length),
		bitset.New(length), //15
	}

	return result
}

func (p *Pattern) translate(i uint, bitSet []*bitset.BitSet) {

	// Define transformation matices.
	t := []int{-p.size / 2, -p.size / 2}
	tInv := []int{p.size / 2, p.size / 2}

	// Transform coordinate.
	x := int(i) / p.size
	y := int(i) % p.size

	for j := 1; j < 8; j++ {

		// Translation
		xt := x + t[0]
		yt := y + t[1]

		// Rotation
		xr := xt
		yr := yt

		switch j {

		case 1:
			xr = p.rot90[0]*xt + p.rot90[1]*yt
			yr = p.rot90[2]*xt + p.rot90[3]*yt

		case 2:
			xr = p.rot180[0]*xt + p.rot180[1]*yt
			yr = p.rot180[2]*xt + p.rot180[3]*yt

		case 3:
			xr = p.rot270[0]*xt + p.rot270[1]*yt
			yr = p.rot270[2]*xt + p.rot270[3]*yt

		case 5:
			xr = p.rot90[0]*xt + p.rot90[1]*yt
			yr = p.rot90[2]*xt + p.rot90[3]*yt

		case 6:
			xr = p.rot180[0]*xt + p.rot180[1]*yt
			yr = p.rot180[2]*xt + p.rot180[3]*yt

		case 7:
			xr = p.rot270[0]*xt + p.rot270[1]*yt
			yr = p.rot270[2]*xt + p.rot270[3]*yt
		}

		// Horizontal mirror.
		xh := xr
		yh := yr

		if j > 3 {
			xh = p.mirror[0]*xr + p.mirror[1]*yr
			yh = p.mirror[2]*xr + p.mirror[3]*yr
		}

		// Translation inverse.
		xt = xh + tInv[0]
		yt = yh + tInv[1]

		// set bit.
		q := uint(xt*p.size + yt)
		bitSet[j].Set(q)
	}
}

func (p *Pattern) combineBitSet(length uint, b1 *bitset.BitSet, b2 *bitset.BitSet) *bitset.BitSet {

	r := bitset.New(length * 2)

	r = r.Union(b2)

	for i := p.nextSetBit(0, b1); i >= 0; i = p.nextSetBit(i+1, b1) {
		// operate on index i here
		r.Set(uint(i) + length)
	}

	return r
}

func (p *Pattern) nextSetBit(index int, b *bitset.BitSet) int {

	r := -1

	for i := index; i < p.size*p.size; i++ {

		if b.Test(uint(i)) {

			r = int(i)

			break
		}
	}

	return r
}

func (p *Pattern) GetHash() int64 {

	b := p.bit2long(p.black)
	w := p.bit2long(p.white)

	r := int64(1)

	for i := 0; i < len(b); i++ {
		r = (p.prime*r + b[i]) % p.celing
	}

	for i := 0; i < len(w); i++ {
		r = (p.prime*r + w[i]) % p.celing
	}

	r = p.prime*r + int64(p.size)

	return r
}

func (p *Pattern) bit2long(b *bitset.BitSet) []int64 {

	l := uint(p.size * p.size)

	// This is to include all bits.
	var s uint
	if l%32 == 0 {
		s = l / 32
	} else {
		s = l/32 + 1
	}

	r := make([]int64, s)

	for i := 0; i < int(s); i++ {

		for j := 0; j < 32; j++ {

			n := i*32 + j

			if uint(n) < l {
				if b.Test(uint(n)) {
					// Max is 4294967295L
					r[i] |= 1 << uint(j)
				}
			}
		}
	}

	return r
}

func (p *Pattern) string(b1 *bitset.BitSet, b2 *bitset.BitSet) {

	var line string

	l := uint(p.size * p.size)

	for i := uint(0); i < l; i++ {

		var c string

		if b1.Test(i) && b2.Test(i) {
			c = "#"
		} else if b1.Test(i) {
			c = "X"
		} else if b2.Test(i) {
			c = "O"
		} else {
			c = "."
		}

		if i%uint(p.size) == 0 && i != 0 {
			fmt.Printf("%v\n", line)
			line = c
		} else {

			line += c
		}
	}
}
