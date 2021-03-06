/*
Package hash provides a library for calculating a canonical hash value of a Go game pattern.

Canonical hash value means same value will be found if the pattern is mirrored, rotated or color changed. Below patterns will have same hash value.

	...X.	...O.	.O...	.....
	..X..	..O..	..O..	.....
	..X.O	..O.X	X.O..	.OO..
	...O.	...X.	.X...	O..X.
	.....	.....	.....	..X..

The edge of the board can be modeled by setting both black and white at the same location.

	.#..O	.....	.....
	.#.X.	.####	####.
	.#.X.	.#...	...#.
	.####	.#XX.	.OO#.
	.....	.#..O	X..#.
*/
package hash

import (
	"github.com/willf/bitset"
)

var (
	mirror = []int{1, 0, 0, -1}
	rot90  = []int{0, -1, 1, 0}
	rot180 = []int{-1, 0, 0, -1}
	rot270 = []int{0, 1, -1, 0}

	maxUint64 uint64 = 18446744073709551615
	prime     uint64 = 13
	wordLimit uint64 = maxUint64 / 2
	celing    uint64 = wordLimit / prime
)

// A Pattern contains data of a Go game situation.
type Pattern struct {
	size int

	black *bitset.BitSet
	white *bitset.BitSet
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
}

func (p *Pattern) SetBlack(x, y int) {

	q := uint(y*p.size + x)

	p.black.Set(q)
}

func (p *Pattern) SetWhite(x, y int) {

	q := uint(y*p.size + x)

	p.white.Set(q)
}

func (p *Pattern) canonical() (black *bitset.BitSet, white *bitset.BitSet) {

	l := uint(p.size * p.size)

	blacks := p.initBitSet(p.black, l)
	whites := p.initBitSet(p.white, l)

	for i := uint(0); i < l; i++ {

		if p.black.Test(i) {
			p.translate(i, blacks)
		}

		if p.white.Test(i) {
			p.translate(i, whites)
		}
	}

	// Reverse color
	for i := 8; i < 16; i++ {

		blacks[i] = whites[i-8]

		whites[i] = blacks[i-8]
	}

	// Find canonical
	bc := p.black
	wc := p.white

	bwc := p.combineBitSet(bc, wc)

	for j := 1; j < 16; j++ {

		a := p.combineBitSet(blacks[j], whites[j])

		b := bwc.Clone()

		// Find different bits.
		b = b.SymmetricDifference(a)

		// Check first true bit.
		n := p.nextSetBit(0, b)
		if n > 0 {

			if bwc.Test(uint(n)) == false {

				bc = blacks[j]
				wc = whites[j]

				bwc = a
			}
		}
	}

	return bc, wc
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
			xr = rot90[0]*xt + rot90[1]*yt
			yr = rot90[2]*xt + rot90[3]*yt

		case 2:
			xr = rot180[0]*xt + rot180[1]*yt
			yr = rot180[2]*xt + rot180[3]*yt

		case 3:
			xr = rot270[0]*xt + rot270[1]*yt
			yr = rot270[2]*xt + rot270[3]*yt

		case 5:
			xr = rot90[0]*xt + rot90[1]*yt
			yr = rot90[2]*xt + rot90[3]*yt

		case 6:
			xr = rot180[0]*xt + rot180[1]*yt
			yr = rot180[2]*xt + rot180[3]*yt

		case 7:
			xr = rot270[0]*xt + rot270[1]*yt
			yr = rot270[2]*xt + rot270[3]*yt
		}

		// Horizontal mirror.
		xh := xr
		yh := yr

		if j > 3 {
			xh = mirror[0]*xr + mirror[1]*yr
			yh = mirror[2]*xr + mirror[3]*yr
		}

		// Translation inverse.
		xt = xh + tInv[0]
		yt = yh + tInv[1]

		// set bit.
		q := uint(xt*p.size + yt)
		bitSet[j].Set(q)
	}
}

func (p *Pattern) combineBitSet(b1 *bitset.BitSet, b2 *bitset.BitSet) *bitset.BitSet {

	return bitset.From(append(b2.Bytes(), b1.Bytes()...))
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

func (p *Pattern) GetHash() uint64 {

	bc, wc := p.canonical()

	b := bc.Bytes()
	w := wc.Bytes()

	r := uint64(1)

	for i := 0; i < len(b); i++ {
		r = hash(r, b[i])
	}

	for i := 0; i < len(w); i++ {
		r = hash(r, w[i])
	}

	r = hash(r, uint64(p.size))

	return r
}

func hash(r, v uint64) uint64 {

	return (prime*r + v%wordLimit) % celing
}

func (p *Pattern) string(b1 *bitset.BitSet, b2 *bitset.BitSet) string {

	var r string

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

		r += c

		if (i+1)%uint(p.size) == 0 && i != 0 {
			r += "\n"
		}
	}

	return r
}
