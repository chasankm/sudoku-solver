package solver

import (
	"fmt"

	"github.com/RoaringBitmap/roaring"
)

type SwordFish struct {
	Up     []*Cell
	Middle []*Cell
	Down   []*Cell
	Mark   *roaring.Bitmap
}

func (s *SwordFish) ColsUnion() []uint32 {
	upIndexes := IndexesBitmap(s.Up)
	middleIndexes := IndexesBitmap(s.Middle)
	downIndexes := IndexesBitmap(s.Down)
	union := ParUnion(upIndexes, middleIndexes, downIndexes)
	return union.ToArray()
}

func (s *SwordFish) Print() {
	fmt.Printf("\nSword Fish\n")
	fmt.Printf("Mark: %s\n", s.Mark.String())
	fmt.Printf("Up Part\n")
	for _, cell := range s.Up {
		fmt.Printf("Cell ID:[%d] --> [%d][%d] Marks: %s\n", cell.ID, cell.Row, cell.Col, cell.Marks.String())
	}
	fmt.Printf("Middle Part\n")
	for _, cell := range s.Middle {
		fmt.Printf("Cell ID:[%d] --> [%d][%d] Marks: %s\n", cell.ID, cell.Row, cell.Col, cell.Marks.String())
	}
	fmt.Printf("Down Part\n")
	for _, cell := range s.Down {
		fmt.Printf("Cell ID:[%d] --> [%d][%d] Marks: %s\n", cell.ID, cell.Row, cell.Col, cell.Marks.String())
	}
}

func (s *SwordFish) Eliminate(b *Board) error {
	for _, index := range s.ColsUnion() {
		col := b.col(int(index))
		for _, c := range col {
			if !c.IsSolved() && !IsCellInCollection(c, s.Up) && !IsCellInCollection(c, s.Middle) && !IsCellInCollection(c, s.Down) {
				c.Marks.AndNot(s.Mark)
				if c.Marks.IsEmpty() {
					return fmt.Errorf("Invalid Board SwordFish: Empty marks: Cell: %+v\n", c)
				}
			}
		}
	}
	return nil
}

func EliminateSwordFish(b *Board) error {
	swordFishes := make([]*SwordFish, 0)
	for i := 0; i < BoardSize; i++ {
		row := b.data[i]
		yes, upCells, mark := HasSwordFishCandidates(row[:])
		if yes {
			yesMiddle, middleCells, index := SearchSwordFishMiddlePart(upCells, mark, i, b)
			if yesMiddle {
				yesDown, swordFish := SearchSwordFishDownPart(upCells, middleCells, mark, index, b)
				if yesDown {
					swordFishes = append(swordFishes, swordFish)
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

func HasSwordFishCandidates(row []*Cell) (bool, []*Cell, *roaring.Bitmap) {
	union := ParUnionCells(UnSolvedCells(row))
	marks := BitmapSingles(union.ToArray())
	for _, mark := range marks {
		yes, cells := IsMarkAppearsTwiceOrThreeInRow(mark, row)
		if yes {
			return true, cells, mark
		}
	}
	return false, nil, nil
}

func IsMarkAppearsTwiceOrThreeInRow(mark *roaring.Bitmap, row []*Cell) (bool, []*Cell) {
	cells := make([]*Cell, 0)
	for _, cell := range row {
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

func SearchSwordFishMiddlePart(upCells []*Cell, mark *roaring.Bitmap, rowIndex int, b *Board) (bool, []*Cell, int) {
	if rowIndex+1 == BoardSize-1 {
		return false, nil, 0
	}
	for i := rowIndex + 1; i < BoardSize; i++ {
		row := b.data[i]
		yes, middleCells := IsMarkAppearsTwiceOrThreeInRow(mark, row[:])
		if yes {
			// Well we found that same mark also only appears twice or three times in this row.
			// Let's also check whether the column indexes also match using sets
			upIndexes := IndexesBitmap(upCells)
			middleIndexes := IndexesBitmap(middleCells)
			intersect := ParIntersect(upIndexes, middleIndexes)
			cardinality := intersect.GetCardinality()
			if cardinality == 1 || cardinality == 2 {
				return true, middleCells, i
			}
		}
	}
	return false, nil, 0
}

func SearchSwordFishDownPart(upCells []*Cell, middleCells []*Cell, mark *roaring.Bitmap, rowIndex int, b *Board) (bool, *SwordFish) {
	if rowIndex+1 == BoardSize-1 {
		return false, nil
	}
	for i := rowIndex + 1; i < BoardSize; i++ {
		row := b.data[i]
		yes, downCells := IsMarkAppearsTwiceOrThreeInRow(mark, row[:])
		if yes {
			upIndexes := IndexesBitmap(upCells)
			middleIndexes := IndexesBitmap(middleCells)
			downIndexes := IndexesBitmap(downCells)
			if ColumnsHasAtLeast2TimesTheValue(upIndexes, middleIndexes, downIndexes) {
				return true, &SwordFish{
					Up:     upCells,
					Middle: middleCells,
					Down:   downCells,
					Mark:   mark,
				}
			}
		}
	}
	return false, nil
}

func IndexesBitmap(cells []*Cell) *roaring.Bitmap {
	indexes := roaring.NewBitmap()
	for _, cell := range cells {
		indexes.Add(uint32(cell.Col))
	}
	return indexes
}

func ColumnsHasAtLeast2TimesTheValue(up *roaring.Bitmap, middle *roaring.Bitmap, down *roaring.Bitmap) bool {
	union := ParUnion(up, middle, down)
	if union.GetCardinality() != 3 {
		return false
	}
	histogram := make(map[uint32]int)
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
