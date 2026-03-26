package solver

import (
	"fmt"
)

type SwordFish struct {
	Up     []*Cell
	Middle []*Cell
	Down   []*Cell
	Mark   CandidateSet
	ByRows bool
}

func (s *SwordFish) TargetIndexes() []int {
	upIndexes := IndexesBitmap(s.Up, s.ByRows)
	middleIndexes := IndexesBitmap(s.Middle, s.ByRows)
	downIndexes := IndexesBitmap(s.Down, s.ByRows)
	union := ParUnion(upIndexes, middleIndexes, downIndexes)
	indexes := union.ToArray()
	for i := range indexes {
		indexes[i]--
	}
	return indexes
}

func (s *SwordFish) Eliminate(b *Board) error {
	for _, index := range s.TargetIndexes() {
		targetUnit := orthogonalLine(b, index, s.ByRows)
		for _, c := range targetUnit {
			if !c.IsSolved() && !IsCellInCollection(c, s.Up) && !IsCellInCollection(c, s.Middle) && !IsCellInCollection(c, s.Down) {
				c.Marks = c.Marks.AndNot(s.Mark)
				if c.Marks.IsEmpty() {
					return fmt.Errorf("invalid board: SwordFish: empty marks: cell: %+v", c)
				}
			}
		}
	}
	return nil
}

func EliminateSwordFish(b *Board) error {
	swordFishes := make([]*SwordFish, 0)
	for _, byRows := range []bool{true, false} {
		for i := 0; i < BoardSize; i++ {
			yes, upCells, mark := HasSwordFishCandidates(line(b, i, byRows))
			if yes {
				yesMiddle, middleCells, index := SearchSwordFishMiddlePart(upCells, mark, i, b, byRows)
				if yesMiddle {
					yesDown, swordFish := SearchSwordFishDownPart(upCells, middleCells, mark, index, b, byRows)
					if yesDown {
						swordFishes = append(swordFishes, swordFish)
					}
				}
			}
		}
	}
	for _, swordFish := range swordFishes {
		if eliminateErr := swordFish.Eliminate(b); eliminateErr != nil {
			return eliminateErr
		}
	}
	return nil
}

func HasSwordFishCandidates(unit []*Cell) (bool, []*Cell, CandidateSet) {
	union := ParUnionCells(UnSolvedCells(unit))
	marks := BitmapSingles(union.ToArray())
	for _, mark := range marks {
		yes, cells := IsMarkAppearsTwiceOrThreeInUnit(mark, unit)
		if yes {
			return true, cells, mark
		}
	}
	return false, nil, 0
}

func IsMarkAppearsTwiceOrThreeInUnit(mark CandidateSet, unit []*Cell) (bool, []*Cell) {
	cells := make([]*Cell, 0)
	for _, cell := range unit {
		if !cell.IsSolved() {
			intersect := ParIntersect(mark, cell.Marks)
			if intersect.GetCardinality() == 1 {
				cells = append(cells, cell)
			}
		}
	}
	if len(cells) == 2 || len(cells) == 3 {
		// Given mark appears only twice or three within the row
		return true, cells
	}
	return false, nil
}

func SearchSwordFishMiddlePart(upCells []*Cell, mark CandidateSet, lineIndex int, b *Board, byRows bool) (bool, []*Cell, int) {
	if lineIndex+1 == BoardSize-1 {
		return false, nil, 0
	}
	for i := lineIndex + 1; i < BoardSize; i++ {
		yes, middleCells := IsMarkAppearsTwiceOrThreeInUnit(mark, line(b, i, byRows))
		if yes {
			upIndexes := IndexesBitmap(upCells, byRows)
			middleIndexes := IndexesBitmap(middleCells, byRows)
			intersect := ParIntersect(upIndexes, middleIndexes)
			cardinality := intersect.GetCardinality()
			if cardinality == 1 || cardinality == 2 {
				return true, middleCells, i
			}
		}
	}
	return false, nil, 0
}

func SearchSwordFishDownPart(upCells []*Cell, middleCells []*Cell, mark CandidateSet, lineIndex int, b *Board, byRows bool) (bool, *SwordFish) {
	if lineIndex+1 == BoardSize-1 {
		return false, nil
	}
	for i := lineIndex + 1; i < BoardSize; i++ {
		yes, downCells := IsMarkAppearsTwiceOrThreeInUnit(mark, line(b, i, byRows))
		if yes {
			upIndexes := IndexesBitmap(upCells, byRows)
			middleIndexes := IndexesBitmap(middleCells, byRows)
			downIndexes := IndexesBitmap(downCells, byRows)
			if ColumnsHasAtLeast2TimesTheValue(upIndexes, middleIndexes, downIndexes) {
				return true, &SwordFish{
					Up:     upCells,
					Middle: middleCells,
					Down:   downCells,
					Mark:   mark,
					ByRows: byRows,
				}
			}
		}
	}
	return false, nil
}

func ColumnsHasAtLeast2TimesTheValue(up CandidateSet, middle CandidateSet, down CandidateSet) bool {
	union := ParUnion(up, middle, down)
	if union.GetCardinality() != 3 {
		return false
	}
	histogram := make(map[int]int)
	// Each column should contain at least 2 shared cells
	for _, index := range union.ToArray() {
		if up.Contains(index) {
			if v, ok := histogram[index]; ok {
				histogram[index] = v + 1
			} else {
				histogram[index] = 1
			}
		}
		if middle.Contains(index) {
			if v, ok := histogram[index]; ok {
				histogram[index] = v + 1
			} else {
				histogram[index] = 1
			}
		}
		if down.Contains(index) {
			if v, ok := histogram[index]; ok {
				histogram[index] = v + 1
			} else {
				histogram[index] = 1
			}
		}
	}
	for _, v := range histogram {
		if v < 2 {
			return false
		}
	}
	return true
}
