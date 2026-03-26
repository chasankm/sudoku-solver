package solver

import (
	"fmt"
)

func EliminateHiddenSingles(units [][]*Cell) error {
	for _, unit := range units {
		unsolved := UnSolvedCells(unit)
		found, single, bitmap := HiddenSingles(unsolved)
		if found {
			single.Marks = single.Marks.And(bitmap)
			if single.Marks.IsEmpty() {
				return fmt.Errorf("invalid board: HS: empty marks: cell: %+v", single)
			}
		}
	}
	return nil
}

func HiddenSingles(unit []*Cell) (bool, *Cell, CandidateSet) {
	for _, single := range unit {
		found, bitmap := IsHiddenSingle(single, unit)
		if found {
			return true, single, bitmap
		}
	}
	return false, nil, 0
}

func IsHiddenSingle(single *Cell, unit []*Cell) (bool, CandidateSet) {
	combinations := BitmapSingles(single.Marks.ToArray())
	for _, t := range combinations {
		if IsCombinationHiddenWithinUnit(t, []*Cell{single}, unit) {
			return true, t
		}
	}
	return false, 0
}

// EliminateHiddenPairs eliminates the marks from the pairs which has exactly and only the same two candidates all
// over the unit. The other candidates could be removed safely from the pairs
func EliminateHiddenPairs(units [][]*Cell) error {
	for _, unit := range units {
		unsolved := UnSolvedCells(unit)
		found, pair, bitmap := HiddenPairs(unsolved)
		if found {
			for _, cell := range pair {
				cell.Marks = cell.Marks.And(bitmap)
				if cell.Marks.IsEmpty() {
					return fmt.Errorf("invalid board: HP: empty marks: cell: %+v", cell)
				}
			}
		}
	}
	return nil
}

// HiddenPairs is a helper method to find any hidden pair over the unit by creating pair combinations and testing
// the pairs whether they are hidden pair or not
func HiddenPairs(unit []*Cell) (bool, []*Cell, CandidateSet) {
	combinations := PairCombinations(unit)
	for _, pairs := range combinations {
		found, bitmap := IsHiddenPair(pairs, unit)
		if found {
			return true, pairs, bitmap
		}
	}
	return false, nil, 0
}

// IsHiddenPair is a helper method to calculate the intersection of the giving pair and looping over the unit
// to get the diff of intersection and other cell marks, and if the final intersection of difference have some
// elements, then the given pair is a hidden pair. Function returns also the marks that has to be kept
func IsHiddenPair(pair []*Cell, unit []*Cell) (bool, CandidateSet) {
	union := ParUnion(pair[0].Marks, pair[1].Marks)
	if union.GetCardinality() <= 2 {
		// No need to elimination, union size already <= 2
		return false, 0
	}
	combinations := BitmapPairs(union.ToArray())
	for _, t := range combinations {
		if IsCombinationHiddenWithinUnit(t, pair, unit) {
			return true, t
		}
	}
	return false, 0
}

// EliminateHiddenTriplets eliminates the marks from the pairs which has exactly and only the same two candidates all
// over the unit. The other candidates could be removed safely from the pairs
func EliminateHiddenTriplets(units [][]*Cell) error {
	for _, unit := range units {
		unsolved := UnSolvedCells(unit)
		found, pair, bitmap := HiddenTriplets(unsolved)
		if found {
			for _, cell := range pair {
				cell.Marks = cell.Marks.And(bitmap)
				if cell.Marks.IsEmpty() {
					return fmt.Errorf("invalid board: HT: empty marks: cell: %+v", cell)
				}
			}
		}
	}
	return nil
}

// HiddenTriplets is a helper method to find any hidden pair over the unit by creating triplet combinations and testing
// the triplets whether they are hidden triplet or not
func HiddenTriplets(unit []*Cell) (bool, []*Cell, CandidateSet) {
	combinations := TripletCombinations(unit)
	for _, pairs := range combinations {
		found, bitmap := IsHiddenTriplet(pairs, unit)
		if found {
			return true, pairs, bitmap
		}
	}
	return false, nil, 0
}

// IsHiddenTriplet is a helper method to calculate the intersection of the giving pair and looping over the unit
// to get the diff of intersection and other cell marks, and if the final intersection of difference have some
// elements, then the given pair is a hidden pair. Function returns also the marks that has to be kept
func IsHiddenTriplet(triplet []*Cell, unit []*Cell) (bool, CandidateSet) {
	union := ParUnion(triplet[0].Marks, triplet[1].Marks, triplet[2].Marks)
	if union.GetCardinality() <= 3 {
		// No need to elimination, union size already <= 3
		return false, 0
	}
	combinations := BitmapTriplets(union.ToArray())
	for _, t := range combinations {
		if IsCombinationHiddenWithinUnit(t, triplet, unit) {
			return true, t
		}
	}
	return false, 0
}

func EliminateHiddenQuads(units [][]*Cell) error {
	for _, unit := range units {
		found, quad, bitmap := HiddenQuads(UnSolvedCells(unit))
		if found {
			for _, cell := range quad {
				cell.Marks = cell.Marks.And(bitmap)
				if cell.Marks.IsEmpty() {
					return fmt.Errorf("invalid board: HQ: empty marks: cell: %+v", cell)
				}
			}
		}
	}
	return nil
}

func HiddenQuads(unit []*Cell) (bool, []*Cell, CandidateSet) {
	combinations := QuadCombinations(unit)
	for _, quads := range combinations {
		found, bitmap := IsHiddenQuad(quads, unit)
		if found {
			return true, quads, bitmap
		}
	}
	return false, nil, 0
}

func IsHiddenQuad(quad []*Cell, unit []*Cell) (bool, CandidateSet) {
	union := ParUnion(quad[0].Marks, quad[1].Marks, quad[2].Marks, quad[3].Marks)
	if union.GetCardinality() <= 4 {
		// No need to elimination, union size already <= 4
		return false, 0
	}
	combinations := BitmapQuads(union.ToArray())
	for _, com := range combinations {
		if IsCombinationHiddenWithinUnit(com, quad, unit) {
			return true, com
		}
	}
	return false, 0
}

func IsCombinationHiddenWithinUnit(bitmap CandidateSet, cells []*Cell, unit []*Cell) bool {
	for _, cell := range unit {
		if !IsCellInCollection(cell, cells) {
			intersect := ParIntersect(bitmap, cell.Marks)
			if !intersect.IsEmpty() {
				return false
			}
		}
	}
	return true
}
