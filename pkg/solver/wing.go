package solver

import (
	"fmt"
)

type XYWing struct {
	Pivot *Cell
	Wings []*Cell
}

func (xy *XYWing) EliminatedMark() (CandidateSet, bool) {
	common := ParIntersect(xy.Wings[0].Marks, xy.Wings[1].Marks)
	if common.GetCardinality() != 1 {
		return 0, false
	}
	if !ParIntersect(common, xy.Pivot.Marks).IsEmpty() {
		return 0, false
	}
	return common, true
}

func (xy *XYWing) Eliminate(b *Board) error {
	marks, ok := xy.EliminatedMark()
	if !ok {
		return nil
	}
	for _, cell := range xy.WingsIntersect(b) {
		cell.Marks = cell.Marks.AndNot(marks)
		if cell.Marks.IsEmpty() {
			return fmt.Errorf("invalid board: XY: empty marks: cell: %+v", cell)
		}
	}
	return nil
}

func (xy *XYWing) Triplet() []*Cell {
	triplet := make([]*Cell, 0, 3)
	triplet = append(triplet, xy.Pivot)
	triplet = append(triplet, xy.Wings...)
	return triplet
}

func (xy *XYWing) WingsIntersect(b *Board) []*Cell {
	intersect := make([]*Cell, 0)
	aUnit := xy.Wings[0].CellUnits(b)
	bUnit := xy.Wings[1].CellUnits(b)
	for _, unit := range aUnit {
		for _, c := range unit {
			if IsCellInCollections(c, bUnit) && c.ID != xy.Wings[0].ID && c.ID != xy.Wings[1].ID && xy.Pivot.ID != c.ID {
				intersect = append(intersect, c)
			}
		}
	}
	return UnSolvedCells(intersect)
}

func EliminateXYWings(unsolved []*Cell, b *Board) error {
	xyWings := make([]*XYWing, 0)
	combinations := TripletCombinations(unsolved)
	for _, triplet := range combinations {
		if IsTripletHasSameCardinality(triplet, 2) && IsTripletUnionHasCardinality(triplet, 3) {
			yes, xyWing := IsTripletXYWingCandidate(triplet, b)
			if yes {
				xyWings = append(xyWings, xyWing)
			}
		}
	}
	for _, xyWing := range xyWings {
		if eliminateErr := xyWing.Eliminate(b); eliminateErr != nil {
			return eliminateErr
		}
	}
	return nil
}

type XYZWing struct {
	Pivot *Cell
	Wings []*Cell
}

func (xyz *XYZWing) Intersect() CandidateSet {
	return ParIntersectCells(xyz.Triplet())
}

func (xyz *XYZWing) Triplet() []*Cell {
	triplet := make([]*Cell, 0, 3)
	triplet = append(triplet, xyz.Pivot)
	triplet = append(triplet, xyz.Wings...)
	return triplet
}

func (xyz *XYZWing) Eliminate(b *Board) error {
	for _, cell := range xyz.XYZIntersect(b) {
		cell.Marks = cell.Marks.AndNot(xyz.Intersect())
		if cell.Marks.IsEmpty() {
			return fmt.Errorf("invalid board: XYZ: empty marks: cell: %+v", cell)
		}
	}
	return nil
}

func (xyz *XYZWing) XYZIntersect(b *Board) []*Cell {
	intersect := make([]*Cell, 0)
	aUnit := xyz.Wings[0].CellUnits(b)
	bUnit := xyz.Wings[1].CellUnits(b)
	cUnit := xyz.Pivot.CellUnits(b)
	for _, unit := range aUnit {
		for _, c := range unit {
			if IsCellInCollections(c, bUnit) && IsCellInCollections(c, cUnit) && (c.ID != xyz.Wings[0].ID && c.ID != xyz.Wings[1].ID && c.ID != xyz.Pivot.ID) {
				intersect = append(intersect, c)
			}
		}
	}
	return UnSolvedCells(intersect)
}

func EliminateXYZWings(unsolved []*Cell, b *Board) error {
	xyzWings := make([]*XYZWing, 0)
	combinations := TripletCombinations(unsolved)
	for _, triplet := range combinations {
		if IsTripletXYZWingBasedOnCardinality(triplet) && IsTripletUnionHasCardinality(triplet, 3) {
			yes, xyzWing := IsTripletXYZWingCandidate(triplet, b)
			if yes {
				xyzWings = append(xyzWings, xyzWing)
			}
		}
	}
	for _, xyzWing := range xyzWings {
		if eliminateErr := xyzWing.Eliminate(b); eliminateErr != nil {
			return eliminateErr
		}
	}
	return nil
}

func IsTripletXYZWingCandidate(triplet []*Cell, b *Board) (bool, *XYZWing) {
	var xyzWing XYZWing
	related := 0
	unRelated := 0
	pairs := PairCombinations(triplet)
	for _, pair := range pairs {
		if !IsPairRelated(pair, b) {
			unRelated++
			xyzWing.Wings = pair
		} else {
			related++
		}
	}
	if related == 2 && unRelated == 1 {
		for _, cell := range triplet {
			if !IsCellInCollection(cell, xyzWing.Wings) {
				xyzWing.Pivot = cell
			}
		}
		// Pivot should have 3 elements and the intersection should be 1
		if xyzWing.Pivot.Marks.GetCardinality() == 3 && xyzWing.Intersect().GetCardinality() == 1 {
			return true, &xyzWing
		}
	}
	return false, nil
}

func IsTripletXYWingCandidate(triplet []*Cell, b *Board) (bool, *XYWing) {
	var xyWing XYWing
	related := 0
	unRelated := 0
	pairs := PairCombinations(triplet)
	for _, pair := range pairs {
		if !IsPairRelated(pair, b) {
			unRelated++
			xyWing.Wings = pair
		} else {
			related++
		}
	}
	if related == 2 && unRelated == 1 {
		for _, cell := range triplet {
			if !IsCellInCollection(cell, xyWing.Wings) {
				xyWing.Pivot = cell
			}
		}
		// Pivot - wings intersect should be 1 element set
		aIntersect := ParIntersect(xyWing.Pivot.Marks, xyWing.Wings[0].Marks)
		bIntersect := ParIntersect(xyWing.Pivot.Marks, xyWing.Wings[1].Marks)
		wingIntersect, ok := xyWing.EliminatedMark()
		if aIntersect.GetCardinality() == 1 &&
			bIntersect.GetCardinality() == 1 &&
			aIntersect != bIntersect &&
			ok &&
			!wingIntersect.IsEmpty() {
			return true, &xyWing
		}
	}
	return false, nil
}

func IsPairRelated(pair []*Cell, board *Board) bool {
	aUnits := pair[0].CellUnits(board)
	for _, unit := range aUnits {
		if IsCellInCollection(pair[1], unit) {
			return true
		}
	}
	return false
}

func IsTripletHasSameCardinality(triplet []*Cell, cardinality int) bool {
	for _, cell := range triplet {
		if cell.Marks.GetCardinality() != cardinality {
			return false
		}
	}
	return true
}

func IsTripletXYZWingBasedOnCardinality(triplet []*Cell) bool {
	// Two cells (wings) might have 2 cardinality while one of them has 3
	cardinality2 := 0
	cardinality3 := 0
	for _, cell := range triplet {
		if cell.Marks.GetCardinality() == 2 {
			cardinality2++
		}
		if cell.Marks.GetCardinality() == 3 {
			cardinality3++
		}
	}
	if cardinality2 == 2 && cardinality3 == 1 {
		return true
	}
	return false
}

func IsTripletUnionHasCardinality(triplet []*Cell, cardinality int) bool {
	union := ParUnionCells(triplet)
	return union.GetCardinality() == cardinality
}

type XWing struct {
	Up     []*Cell
	Down   []*Cell
	Mark   CandidateSet
	ByRows bool
}

func line(b *Board, index int, byRows bool) []*Cell {
	if byRows {
		return b.row(index)
	}
	return b.col(index)
}

func orthogonalLine(b *Board, index int, byRows bool) []*Cell {
	if byRows {
		return b.col(index)
	}
	return b.row(index)
}

func IndexesBitmap(cells []*Cell, byRows bool) CandidateSet {
	var indexes CandidateSet
	for _, cell := range cells {
		if byRows {
			indexes = indexes.Add(cell.Col + 1)
		} else {
			indexes = indexes.Add(cell.Row + 1)
		}
	}
	return indexes
}

func (x *XWing) TargetUnit(b *Board, wingIndex int) ([]*Cell, error) {
	if x.ByRows {
		if x.Up[wingIndex].Col != x.Down[wingIndex].Col {
			return nil, fmt.Errorf("inconsistent X-Wing col indexes: %+v", x)
		}
		return orthogonalLine(b, x.Up[wingIndex].Col, x.ByRows), nil
	}
	if x.Up[wingIndex].Row != x.Down[wingIndex].Row {
		return nil, fmt.Errorf("inconsistent X-Wing row indexes: %+v", x)
	}
	return orthogonalLine(b, x.Up[wingIndex].Row, x.ByRows), nil
}

func (x *XWing) Eliminate(b *Board) error {
	for wingIndex := range []int{0, 1} {
		targetUnit, targetUnitErr := x.TargetUnit(b, wingIndex)
		if targetUnitErr != nil {
			return targetUnitErr
		}
		for _, cell := range targetUnit {
			if !cell.IsSolved() && cell.ID != x.Up[wingIndex].ID && cell.ID != x.Down[wingIndex].ID {
				cell.Marks = cell.Marks.AndNot(x.Mark)
				if cell.Marks.IsEmpty() {
					return fmt.Errorf("invalid board: X: empty marks: cell: %+v", cell)
				}
			}
		}
	}
	return nil
}

func EliminateXWings(b *Board) error {
	xWings := make([]*XWing, 0)
	for _, byRows := range []bool{true, false} {
		for i := 0; i < BoardSize; i++ {
			yes, cells, mark := HasXCandidates(line(b, i, byRows))
			if yes {
				y, xWing := SearchDownPart(cells, mark, i, b, byRows)
				if y {
					xWings = append(xWings, xWing)
				}
			}
		}
	}
	for _, xWing := range xWings {
		if eliminateErr := xWing.Eliminate(b); eliminateErr != nil {
			return eliminateErr
		}
	}
	return nil
}

func HasXCandidates(unit []*Cell) (bool, []*Cell, CandidateSet) {
	union := ParUnionCells(UnSolvedCells(unit))
	marks := BitmapSingles(union.ToArray())
	for _, mark := range marks {
		yes, cells := IsMarkAppearsTwiceInUnit(mark, unit)
		if yes {
			return true, cells, mark
		}
	}
	return false, nil, 0
}

func IsMarkAppearsTwiceInUnit(mark CandidateSet, unit []*Cell) (bool, []*Cell) {
	cells := make([]*Cell, 0)
	for _, cell := range unit {
		if !cell.IsSolved() {
			intersect := ParIntersect(mark, cell.Marks)
			if intersect.GetCardinality() == 1 {
				cells = append(cells, cell)
			}
		}
	}
	if len(cells) == 2 {
		// Given mark appears only twice within the row
		return true, cells
	}
	return false, nil
}

func SearchDownPart(upCells []*Cell, mark CandidateSet, lineIndex int, b *Board, byRows bool) (bool, *XWing) {
	if lineIndex+1 == BoardSize-1 {
		return false, nil
	}
	for i := lineIndex + 1; i < BoardSize; i++ {
		yes, downCells := IsMarkAppearsTwiceInUnit(mark, line(b, i, byRows))
		if yes {
			upIndexes := IndexesBitmap(upCells, byRows)
			downIndexes := IndexesBitmap(downCells, byRows)
			intersect := ParIntersect(upIndexes, downIndexes)
			if intersect.GetCardinality() == 2 {
				return true, &XWing{
					Up:     upCells,
					Down:   downCells,
					Mark:   mark,
					ByRows: byRows,
				}
			}
		}
	}
	return false, nil
}
