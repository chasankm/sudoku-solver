package solver

import (
	"github.com/RoaringBitmap/roaring"
)

// PairCombinations simply creates all unique ordered pair combinations of given cell unit
func PairCombinations(input []*Cell) [][]*Cell {
	pairs := make([][]*Cell, 0)
	for i := 0; i < len(input) - 1; i++ {
		for j := i + 1; j < len(input); j++ {
			pair := make([]*Cell, 0, 2)
			pair = append(pair, input[i], input[j])
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

// TripletCombinations simply creates all unique ordered triple combinations of given cells unit
func TripletCombinations(input []*Cell) [][]*Cell {
	triplets := make([][]*Cell, 0)
	for i := 0; i < len(input) - 2; i++ {
		for j := i + 1; j < len(input) - 1; j++ {
			for k := j + 1; k < len(input); k++ {
				triplet := make([]*Cell, 0, 3)
				triplet = append(triplet, input[i], input[j], input[k])
				triplets = append(triplets, triplet)
			}
		}
	}
	return triplets
}

// QuadCombinations simply creates all unique ordered quad combinations of given cells unit
func QuadCombinations(input []*Cell) [][]*Cell {
	quads := make([][]*Cell, 0)
	for i := 0; i < len(input) - 3; i++ {
		for j := i + 1; j < len(input) - 2; j++ {
			for k := j + 1; k < len(input) - 1; k++ {
				for l := k + 1; l < len(input); l++ {
					quad := make([]*Cell, 0, 4)
					quad = append(quad, input[i], input[j], input[k], input[l])
					quads = append(quads, quad)
				}
			}
		}
	}
	return quads
}

// BitmapSingles simply creates all unique ordered single combinations in given uint32 array as roaring.Bitmap instance
func BitmapSingles(in []uint32) []*roaring.Bitmap {
	singles := make([]*roaring.Bitmap, 0, len(in))
	for _, v := range in {
		singles = append(singles, roaring.BitmapOf(v))
	}
	return singles
}

// BitmapPairs simply creates all unique ordered pair combinations in given uint32 array as roaring.Bitmap instance
func BitmapPairs(in []uint32) []*roaring.Bitmap {
	pairs := make([]*roaring.Bitmap, 0)
	for i := 0; i < len(in) - 1; i++ {
		for j := i + 1; j < len(in); j++ {
			pairs = append(pairs, roaring.BitmapOf(in[i], in[j]))
		}
	}
	return pairs
}

// BitmapTriplets simply creates all unique ordered triplet combinations in given uint32 array as roaring.Bitmap instance
func BitmapTriplets(in []uint32) []*roaring.Bitmap {
	triplets := make([]*roaring.Bitmap, 0)
	for i := 0; i < len(in) - 2; i++ {
		for j := i + 1; j < len(in) - 1; j++ {
			for k := j + 1; k < len(in); k++ {
				triplets = append(triplets, roaring.BitmapOf(in[i], in[j], in[k]))
			}
		}
	}
	return triplets
}

// BitmapQuads simply creates all unique ordered quad combinations in given uint32 array as roaring.Bitmap instance
func BitmapQuads(in []uint32) []*roaring.Bitmap {
	quads := make([]*roaring.Bitmap, 0)
	for i := 0; i < len(in) - 3; i++ {
		for j := i + 1; j < len(in) - 2; j++ {
			for k := j + 1; k < len(in) - 1; k++ {
				for l := k + 1; l < len(in); l++ {
					quads = append(quads, roaring.BitmapOf(in[i], in[j], in[k], in[l]))
				}
			}
		}
	}
	return quads
}

