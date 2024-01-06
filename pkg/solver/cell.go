package solver

import "github.com/RoaringBitmap/roaring"

// Value is the value type of the solved cell which is simply a byte
type Value byte

// Cell is the struct of cell keeping unique id, row and col ids solved value and current marks/candidates
type Cell struct {
	ID    int
	Row   int
	Col   int
	Value Value
	Marks *roaring.Bitmap
}

// CellUnits returns the related cells row, col and box
func (c *Cell) CellUnits(b *Board) [][]*Cell {
	cells := make([][]*Cell, 0)
	cells = append(cells, b.row(c.Row), b.col(c.Col), b.box(c.Row, c.Col))
	return cells
}

// ComputeCellMarks computes the candidates/marks of the current cell
func (c *Cell) ComputeCellMarks(b *Board) *roaring.Bitmap {
	row := FindMarksOfUnit(b.row(c.Row))
	col := FindMarksOfUnit(b.col(c.Col))
	box := FindMarksOfUnit(b.box(c.Row, c.Col))
	return ParIntersect(row, col, box)
}

// IsValid checks the validity of given v value to put in cell
func (c *Cell) IsValid(b *Board, v Value) bool {
	for _, r := range b.row(c.Row) {
		if r.Value == v {
			return false
		}
	}
	for _, r := range b.col(c.Col) {
		if r.Value == v {
			return false
		}
	}
	for _, r := range b.box(c.Row, c.Col) {
		if r.Value == v {
			return false
		}
	}
	return true
}

// IsSolved simply returns whether the cell is already solved or not
func (c *Cell) IsSolved() bool {
	return c.Value != EmptyCellValue
}

// MarksLength simply returns the length of the current marks bitmap
func (c *Cell) MarksLength() int {
	return len(c.Marks.ToArray())
}
