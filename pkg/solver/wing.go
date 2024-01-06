package solver

import (
	"fmt"

	"github.com/RoaringBitmap/roaring"
)

type XYWing struct {
	Pivot *Cell
	Wings []*Cell
}

func (xy *XYWing) Union() *roaring.Bitmap {
	return ParUnionCells(xy.Triplet())
}

func (xy *XYWing) Eliminate(b *Board) error {
	union := xy.Union()
	union.AndNot(xy.Pivot.Marks)
	marks := union.Clone()
	for _, cell := range xy.WingsIntersect(b) {
		cell.Marks.AndNot(marks)
		if cell.Marks.IsEmpty() {
			return fmt.Errorf("Invalid Board XY: Empty marks: Cell: %+v\n", cell)
		}
	}
	return nil
}

func (xy *XYWing) Triplet() []*Cell {
	triplet := make([]*Cell, 0, 3)
	triplet = append(triplet, xy.Pivot)
	for _, wing := range xy.Wings {
		triplet = append(triplet, wing)
	}
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

func (xy *XYWing) Print() {
	fmt.Printf("\nXY Wing\n")
	fmt.Printf("Pivot: %+v\n", xy.Pivot)
	fmt.Printf("Wings\n")
	for _, wing := range xy.Wings {
		fmt.Printf("Wing: %+v\n", wing)
	}
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

func (xyz *XYZWing) Intersect() *roaring.Bitmap {
	return ParIntersectCells(xyz.Triplet())
}

func (xyz *XYZWing) Triplet() []*Cell {
	triplet := make([]*Cell, 0, 3)
	triplet = append(triplet, xyz.Pivot)
	for _, wing := range xyz.Wings {
		triplet = append(triplet, wing)
	}
	return triplet
}

func (xyz *XYZWing) Eliminate(b *Board) error {
	for _, cell := range xyz.XYZIntersect(b) {
		cell.Marks.AndNot(xyz.Intersect())
		if cell.Marks.IsEmpty() {
			return fmt.Errorf("Invalid Board XYZ: Empty marks: Cell: %+v\n", cell)
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

func (xyz *XYZWing) Print() {
	fmt.Printf("XYZ Wing\n")
	fmt.Printf("Pivot: %+v\n", xyz.Pivot)
	fmt.Printf("Wings\n")
	for _, wing := range xyz.Wings {
		fmt.Printf("Wing: %+v\n", wing)
	}
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
		if aIntersect.GetCardinality() == 1 && bIntersect.GetCardinality() == 1 {
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
		if cell.Marks.GetCardinality() != uint64(cardinality) {
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
	return union.GetCardinality() == uint64(cardinality)
}

type XWing struct {
	Up   []*Cell
	Down []*Cell
	Mark *roaring.Bitmap
}

func (x *XWing) LeftCol(b *Board) ([]*Cell, error) {
	if x.Up[0].Col != x.Down[0].Col {
		return nil, fmt.Errorf("Inconsistency on LEFT col indexes: %+v\n", x)
	}
	return b.col(x.Up[0].Col), nil
}

func (x *XWing) RightCol(b *Board) ([]*Cell, error) {
	if x.Up[1].Col != x.Down[1].Col {
		return nil, fmt.Errorf("Inconsistency on RIGHT col indexes: %+v\n", x)
	}
	return b.col(x.Up[1].Col), nil
}

func (x *XWing) Eliminate(b *Board) error {
	leftColumn, leftColErr := x.LeftCol(b)
	if leftColErr != nil {
		return leftColErr
	}
	for _, cell := range leftColumn {
		if !cell.IsSolved() && cell.ID != x.Up[0].ID && cell.ID != x.Down[0].ID {
			cell.Marks.AndNot(x.Mark)
			if cell.Marks.IsEmpty() {
				return fmt.Errorf("Invalid Board X (LEFT): Empty marks: Cell: %+v\n", cell)
			}
		}
	}

	rightColumn, rightColErr := x.RightCol(b)
	if rightColErr != nil {
		return rightColErr
	}
	for _, cell := range rightColumn {
		if !cell.IsSolved() && cell.ID != x.Up[1].ID && cell.ID != x.Down[1].ID {
			cell.Marks.AndNot(x.Mark)
			if cell.Marks.IsEmpty() {
				return fmt.Errorf("Invalid Board X (RIGHT): Empty marks: Cell: %+v\n", cell)
			}
		}
	}
	return nil
}

func (x *XWing) Print() {
	fmt.Printf("\nX Wing\n")
	fmt.Printf("Mark: %s\n", x.Mark.String())
	fmt.Printf("Up Part\n")
	for _, cell := range x.Up {
		fmt.Printf("Cell [%d][%d] Marks: %s\n", cell.Row, cell.Col, cell.Marks.String())
	}
	fmt.Printf("Down Part\n")
	for _, cell := range x.Down {
		fmt.Printf("Cell [%d][%d] Marks: %s\n", cell.Row, cell.Col, cell.Marks.String())
	}
}

func EliminateXWings(b *Board) error {
	xWings := make([]*XWing, 0)
	for i := 0; i < BoardSize; i++ {
		row := b.data[i]
		yes, cells, mark := HasXCandidates(row[:])
		if yes {
			y, xWing := SearchDownPart(cells, mark, i, b)
			if y {
				xWings = append(xWings, xWing)
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

func HasXCandidates(row []*Cell) (bool, []*Cell, *roaring.Bitmap) {
	union := ParUnionCells(UnSolvedCells(row))
	marks := BitmapSingles(union.ToArray())
	for _, mark := range marks {
		yes, cells := IsMarkAppearsTwiceInRow(mark, row)
		if yes {
			return true, cells, mark
		}
	}
	return false, nil, nil
}

func IsMarkAppearsTwiceInRow(mark *roaring.Bitmap, row []*Cell) (bool, []*Cell) {
	cells := make([]*Cell, 0)
	for _, cell := range row {
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

func SearchDownPart(upCells []*Cell, mark *roaring.Bitmap, rowIndex int, b *Board) (bool, *XWing) {
	if rowIndex+1 == BoardSize-1 {
		return false, nil
	}
	for i := rowIndex + 1; i < BoardSize; i++ {
		row := b.data[i]
		yes, downCells := IsMarkAppearsTwiceInRow(mark, row[:])
		if yes {
			// Well we found that same mark also only appears twice in this row.
			// Let's also check whether the column indexes also match using sets
			upIndexes := IndexesBitmap(upCells)
			downIndexes := IndexesBitmap(downCells)
			intersect := ParIntersect(upIndexes, downIndexes)
			if intersect.GetCardinality() == 2 {
				// Indexes also match perfectly
				return true, &XWing{
					Up:   upCells,
					Down: downCells,
					Mark: mark,
				}
			}
		}
	}
	return false, nil
}
