package solver

import "fmt"

// EliminateNakedPairs creates pairs combinations of marks/candidates on each unit (row, col or box)
// if there are any pairs having the cardinality 2 (Unions of the three sets has exactly 2 different elements)
// Then the other cell candidates within the unit having one of these elements is safely eliminated
// Returns the number of eliminated candidates
func EliminateNakedPairs(units [][]*Cell) error {
	for _, unit := range units {
		unsolved := UnSolvedCells(unit)
		pairs := PairCombinations(unsolved)
		found, pair := IsNakedPairs(pairs)
		if found {
			marks := ParUnionCells(pair)
			for _, cell := range unit {
				if !IsCellInCollection(cell, pair) && !cell.IsSolved() {
					cell.Marks.AndNot(marks)
					if cell.Marks.IsEmpty() {
						return fmt.Errorf("Invalid Board NP: Empty marks: Cell: %+v\n", cell)
					}
				}
			}
		}
	}
	return nil
}

// IsNakedPairs simply checks whether there are any naked pairs or not in the given combination
func IsNakedPairs(combinations [][]*Cell) (bool, []*Cell) {
	for _, pair := range combinations {
		if pair[0].Marks.GetCardinality() == 2 && pair[1].Marks.GetCardinality() == 2 {
			union := ParUnionCells(pair)
			if union.GetCardinality() == 2 {
				return true, pair
			}
		}
	}
	return false, nil
}

// EliminateNakedTriplets creates triple combinations of marks/candidates on each unit (row, col or box)
// if there are any triplets having the cardinality 3 (Unions of the three sets has exactly 3 different elements)
// Then the other cell candidates within the unit having one of these elements is safely eliminated
func EliminateNakedTriplets(units [][]*Cell) error {
	for _, unit := range units {
		unsolved := UnSolvedCells(unit)
		triplets := TripletCombinations(unsolved)
		found, triplet := IsNakedTriplet(triplets)
		if found {
			marks := ParUnionCells(triplet)
			for _, cell := range unit {
				if !IsCellInCollection(cell, triplet) && !cell.IsSolved() {
					cell.Marks.AndNot(marks)
					if cell.Marks.IsEmpty() {
						return fmt.Errorf("Invalid Board NT: Empty marks: Cell: %+v\n", cell)
					}
				}
			}
		}
	}
	return nil
}

// IsNakedTriplet simply checks whether there are any naked triplets or not in the given combination
func IsNakedTriplet(combinations [][]*Cell) (bool, []*Cell) {
	for _, triplet := range combinations {
		union := ParUnionCells(triplet)
		if union.GetCardinality() == 3 {
			return true, triplet
		}
	}
	return false, nil
}

// EliminateNakedQuads creates quad combinations of marks/candidates on each unit (row, col or box)
// if there are any quads having the cardinality 4 (Unions of the three sets has exactly 4 different elements)
// Then the other cell candidates within the unit having one of these elements is safely eliminated
// Returns the number of eliminated candidates
func EliminateNakedQuads(units [][]*Cell) error {
	for _, unit := range units {
		unsolved := UnSolvedCells(unit)
		quads := QuadCombinations(unsolved)
		found, quad := IsNakedQuad(quads)
		if found {
			marks := ParUnionCells(quad)
			for _, cell := range unit {
				if !IsCellInCollection(cell, quad) && !cell.IsSolved() {
					cell.Marks.AndNot(marks)
					if cell.Marks.IsEmpty() {
						return fmt.Errorf("Invalid Board NQ: Empty marks: Cell: %+v\n", cell)
					}
				}
			}
		}
	}
	return nil
}

// IsNakedQuad simply checks whether there are any naked quads or not in the given combination
func IsNakedQuad(combinations [][]*Cell) (bool, []*Cell) {
	for _, quad := range combinations {
		union := ParUnionCells(quad)
		if union.GetCardinality() == 4 {
			return true, quad
		}
	}
	return false, nil
}
