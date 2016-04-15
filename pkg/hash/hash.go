/*
Package hash provides a library for calculating a canonical hash value for a Go game pattern.

Canonical hash value means same value will be found if the pattern is mirrored, rotated or color changed. Below patterns will have same hash value.

	...X.	...O.	.O...	.....
	..X..	..O..	..O..	.....
	..X.O	..O.X	X.O..	.OO..
	...O.	...X.	.X...	O..X.
	.....	.....	.....	..X..
*/
package hash

import (
	"github.com/willf/bitset"
)

type Pattern struct {
	bounding_box_size int

	Black_pattern *bitset.BitSet
	White_pattern *bitset.BitSet

	black_BitSets []*bitset.BitSet
	white_BitSets []*bitset.BitSet

	horizontal  []int
	rotation90  []int
	rotation180 []int
	rotation270 []int

	prime  int64
	celing int64
}

func (pattern *Pattern) Initialize(size int) {

	pattern.bounding_box_size = size

	pattern.Black_pattern = bitset.New(uint(size * size))

	pattern.White_pattern = bitset.New(uint(size * size))

	pattern.horizontal = []int{1, 0, 0, -1}
	pattern.rotation90 = []int{0, -1, 1, 0}
	pattern.rotation180 = []int{-1, 0, 0, -1}
	pattern.rotation270 = []int{0, 1, -1, 0}

	pattern.prime = 13
	pattern.celing = (9223372036854775807 - 4294967295) / pattern.prime
}

func (pattern *Pattern) SetBlack(x, y int) {

	q := uint(y*pattern.bounding_box_size + x)

	pattern.Black_pattern.Set(q)
}

func (pattern *Pattern) SetWhite(x, y int) {

	q := uint(y*pattern.bounding_box_size + x)

	pattern.White_pattern.Set(q)
}

func (pattern *Pattern) Canonical() {

	// Find canonical pattern.
	length := uint(pattern.bounding_box_size * pattern.bounding_box_size)

	// Initialize BitSet arrays.
	pattern.black_BitSets = pattern.initializeBitSetArrays(pattern.Black_pattern, length)
	pattern.white_BitSets = pattern.initializeBitSetArrays(pattern.White_pattern, length)

	for i := uint(0); i < length; i++ {

		if pattern.Black_pattern.Test(i) {
			pattern.translateBit(i, pattern.black_BitSets)
		}

		if pattern.White_pattern.Test(i) {
			pattern.translateBit(i, pattern.white_BitSets)
		}
	}

	// Reverse color
	for i := 8; i < 16; i++ {

		pattern.black_BitSets[i] = pattern.white_BitSets[i-8]

		pattern.white_BitSets[i] = pattern.black_BitSets[i-8]
	}

	// Find canonical
	black_canonical := pattern.Black_pattern
	white_canonical := pattern.White_pattern

	BW_canonical := pattern.combineBitSets(length, black_canonical, white_canonical)

	for j := 1; j < 16; j++ {

		BW_current := pattern.combineBitSets(length, pattern.black_BitSets[j], pattern.white_BitSets[j])

		BW_compare := bitset.New(length * 2)
		BW_compare = BW_compare.Union(BW_canonical)

		BW_compare = BW_compare.SymmetricDifference(BW_current) // Find different bits.

		// Check first true bit.
		firstDiffenrence := pattern.nextSetBit(0, BW_compare)
		if firstDiffenrence >= 0 {

			if BW_canonical.Test(uint(firstDiffenrence)) == false {

				black_canonical = pattern.black_BitSets[j]
				white_canonical = pattern.white_BitSets[j]

				BW_canonical = BW_current
			}
		}
	}

	// Set canonical pattern
	pattern.Black_pattern = black_canonical
	pattern.White_pattern = white_canonical
}

func (pattern *Pattern) initializeBitSetArrays(bitSet *bitset.BitSet, length uint) []*bitset.BitSet {

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

func (pattern *Pattern) translateBit(i uint, bitSet []*bitset.BitSet) {

	// Define transformation matices.
	tranlation := []int{-pattern.bounding_box_size / 2, -pattern.bounding_box_size / 2}
	tranlationInverse := []int{pattern.bounding_box_size / 2, pattern.bounding_box_size / 2}

	// Transform coordinate.
	x := int(i) / pattern.bounding_box_size
	y := int(i) % pattern.bounding_box_size

	for j := 1; j < 8; j++ {

		// Translation
		xt := x + tranlation[0]
		yt := y + tranlation[1]

		// Rotation
		xr := xt
		yr := yt

		switch j {

		case 1:
			xr = pattern.rotation90[0]*xt + pattern.rotation90[1]*yt
			yr = pattern.rotation90[2]*xt + pattern.rotation90[3]*yt

		case 2:
			xr = pattern.rotation180[0]*xt + pattern.rotation180[1]*yt
			yr = pattern.rotation180[2]*xt + pattern.rotation180[3]*yt

		case 3:
			xr = pattern.rotation270[0]*xt + pattern.rotation270[1]*yt
			yr = pattern.rotation270[2]*xt + pattern.rotation270[3]*yt

		case 5:
			xr = pattern.rotation90[0]*xt + pattern.rotation90[1]*yt
			yr = pattern.rotation90[2]*xt + pattern.rotation90[3]*yt

		case 6:
			xr = pattern.rotation180[0]*xt + pattern.rotation180[1]*yt
			yr = pattern.rotation180[2]*xt + pattern.rotation180[3]*yt

		case 7:
			xr = pattern.rotation270[0]*xt + pattern.rotation270[1]*yt
			yr = pattern.rotation270[2]*xt + pattern.rotation270[3]*yt
		}

		// Horizontal mirror.
		xh := xr
		yh := yr

		if j > 3 {
			xh = pattern.horizontal[0]*xr + pattern.horizontal[1]*yr
			yh = pattern.horizontal[2]*xr + pattern.horizontal[3]*yr
		}

		// Translation inverse.
		xt = xh + tranlationInverse[0]
		yt = yh + tranlationInverse[1]

		// set bit.
		q := uint(xt*pattern.bounding_box_size + yt)
		bitSet[j].Set(q)
	}
}

func (pattern *Pattern) combineBitSets(length uint, bitSet1 *bitset.BitSet, bitSet2 *bitset.BitSet) *bitset.BitSet {

	result := bitset.New(length * 2)

	result = result.Union(bitSet2)

	for i := pattern.nextSetBit(0, bitSet1); i >= 0; i = pattern.nextSetBit(i+1, bitSet1) {
		// operate on index i here
		result.Set(uint(i) + length)
	}

	return result
}

func (pattern *Pattern) nextSetBit(index int, bitSet *bitset.BitSet) int {

	result := -1

	for i := index; i < pattern.bounding_box_size*pattern.bounding_box_size; i++ {

		if bitSet.Test(uint(i)) {

			result = int(i)

			break
		}
	}

	return result
}

func (pattern *Pattern) GetHashMax() int64 {

	blacks := pattern.bits2longs(pattern.Black_pattern)
	whites := pattern.bits2longs(pattern.White_pattern)

	result := int64(1)

	for i := 0; i < len(blacks); i++ {
		result = (pattern.prime*result + blacks[i]) % pattern.celing
	}

	for i := 0; i < len(whites); i++ {
		result = (pattern.prime*result + whites[i]) % pattern.celing
	}

	result = pattern.prime*result + int64(pattern.bounding_box_size)

	return result
}

func (pattern *Pattern) bits2longs(bs *bitset.BitSet) []int64 {

	length := uint(pattern.bounding_box_size * pattern.bounding_box_size)

	// This is to include all bits.
	var arSize uint
	if length%32 == 0 {
		arSize = length / 32
	} else {
		arSize = length/32 + 1
	}

	result := make([]int64, arSize)

	for i := 0; i < int(arSize); i++ {

		for j := 0; j < 32; j++ {

			index := i*32 + j

			if uint(index) < length {
				if bs.Test(uint(index)) {
					result[i] |= 1 << uint(j) // Max is 4294967295L
				}
			}
		}
	}

	return result
}

func (pattern *Pattern) LogBoard(bitSet1 *bitset.BitSet, bitSet2 *bitset.BitSet) {

	var line string

	l := uint(pattern.bounding_box_size * pattern.bounding_box_size)

	for i := uint(0); i < l; i++ {

		var c string

		if bitSet1.Test(i) && bitSet2.Test(i) {
			c = "#"
		} else if bitSet1.Test(i) {
			c = "X"
		} else if bitSet2.Test(i) {
			c = "O"
		} else {
			c = "."
		}

		if i%uint(pattern.bounding_box_size) == 0 && i != 0 {
			//slog.Debug("%v", line)
			line = c
		} else {

			line += c
		}
	}
}
