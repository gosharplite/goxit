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
