package solver

import "fmt"

func EliminateLockedCandidates(b *Board) error {
	for digit := 1; digit <= BoardSize; digit++ {
		mark := CandidateSetOf(digit)
		if err := eliminatePointingLockedCandidates(b, mark); err != nil {
			return err
		}
		if err := eliminateClaimingLockedCandidates(b, mark); err != nil {
			return err
		}
	}
	return nil
}

func eliminatePointingLockedCandidates(b *Board, mark CandidateSet) error {
	for boxRow := 0; boxRow < BoardSize; boxRow += BlockSize {
		for boxCol := 0; boxCol < BoardSize; boxCol += BlockSize {
			cells := candidateCellsForMark(b.box(boxRow, boxCol), mark)
			if len(cells) < 2 {
				continue
			}

			if sameRow, row := confinedRow(cells); sameRow {
				if err := eliminateMarkFromRowOutsideBox(b, row, boxCol/BlockSize, cells, mark); err != nil {
					return err
				}
			}
			if sameCol, col := confinedCol(cells); sameCol {
				if err := eliminateMarkFromColOutsideBox(b, col, boxRow/BlockSize, cells, mark); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func eliminateClaimingLockedCandidates(b *Board, mark CandidateSet) error {
	for row := 0; row < BoardSize; row++ {
		cells := candidateCellsForMark(b.row(row), mark)
		if len(cells) < 2 {
			continue
		}
		if sameBox, box := confinedBox(cells); sameBox {
			if err := eliminateMarkFromBoxOutsideRow(b, row, box, cells, mark); err != nil {
				return err
			}
		}
	}

	for col := 0; col < BoardSize; col++ {
		cells := candidateCellsForMark(b.col(col), mark)
		if len(cells) < 2 {
			continue
		}
		if sameBox, box := confinedBox(cells); sameBox {
			if err := eliminateMarkFromBoxOutsideCol(b, col, box, cells, mark); err != nil {
				return err
			}
		}
	}

	return nil
}

func candidateCellsForMark(cells []*Cell, mark CandidateSet) []*Cell {
	candidates := make([]*Cell, 0)
	for _, cell := range cells {
		if !cell.IsSolved() && !ParIntersect(cell.Marks, mark).IsEmpty() {
			candidates = append(candidates, cell)
		}
	}
	return candidates
}

func confinedRow(cells []*Cell) (bool, int) {
	if len(cells) == 0 {
		return false, -1
	}
	row := cells[0].Row
	for _, cell := range cells[1:] {
		if cell.Row != row {
			return false, -1
		}
	}
	return true, row
}

func confinedCol(cells []*Cell) (bool, int) {
	if len(cells) == 0 {
		return false, -1
	}
	col := cells[0].Col
	for _, cell := range cells[1:] {
		if cell.Col != col {
			return false, -1
		}
	}
	return true, col
}

func confinedBox(cells []*Cell) (bool, int) {
	if len(cells) == 0 {
		return false, -1
	}
	box := boxIndex(cells[0].Row, cells[0].Col)
	for _, cell := range cells[1:] {
		if boxIndex(cell.Row, cell.Col) != box {
			return false, -1
		}
	}
	return true, box
}

func eliminateMarkFromRowOutsideBox(b *Board, row int, box int, protected []*Cell, mark CandidateSet) error {
	for _, cell := range b.row(row) {
		if boxIndex(cell.Row, cell.Col) == box || IsCellInCollection(cell, protected) || cell.IsSolved() {
			continue
		}
		if err := eliminateMarkFromCell(cell, mark, "Locked Candidates"); err != nil {
			return err
		}
	}
	return nil
}

func eliminateMarkFromColOutsideBox(b *Board, col int, box int, protected []*Cell, mark CandidateSet) error {
	for _, cell := range b.col(col) {
		if boxIndex(cell.Row, cell.Col) == box || IsCellInCollection(cell, protected) || cell.IsSolved() {
			continue
		}
		if err := eliminateMarkFromCell(cell, mark, "Locked Candidates"); err != nil {
			return err
		}
	}
	return nil
}

func eliminateMarkFromBoxOutsideRow(b *Board, row int, box int, protected []*Cell, mark CandidateSet) error {
	boxRow := (box / BlockSize) * BlockSize
	boxCol := (box % BlockSize) * BlockSize
	for _, cell := range b.box(boxRow, boxCol) {
		if cell.Row == row || IsCellInCollection(cell, protected) || cell.IsSolved() {
			continue
		}
		if err := eliminateMarkFromCell(cell, mark, "Locked Candidates"); err != nil {
			return err
		}
	}
	return nil
}

func eliminateMarkFromBoxOutsideCol(b *Board, col int, box int, protected []*Cell, mark CandidateSet) error {
	boxRow := (box / BlockSize) * BlockSize
	boxCol := (box % BlockSize) * BlockSize
	for _, cell := range b.box(boxRow, boxCol) {
		if cell.Col == col || IsCellInCollection(cell, protected) || cell.IsSolved() {
			continue
		}
		if err := eliminateMarkFromCell(cell, mark, "Locked Candidates"); err != nil {
			return err
		}
	}
	return nil
}

func eliminateMarkFromCell(cell *Cell, mark CandidateSet, strategy string) error {
	if ParIntersect(cell.Marks, mark).IsEmpty() {
		return nil
	}
	cell.Marks = cell.Marks.AndNot(mark)
	if cell.Marks.IsEmpty() {
		return fmt.Errorf("invalid board: %s: empty marks: cell: %+v", strategy, cell)
	}
	return nil
}

func boxIndex(row int, col int) int {
	return (row/BlockSize)*BlockSize + (col / BlockSize)
}
