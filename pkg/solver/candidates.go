package solver

import "fmt"

func (b *Board) initializeCandidates() error {
	return b.computeAllMarks()
}

func (b *Board) resolveSingles() (bool, error) {
	changed := false

	for {
		for i := 0; i < BoardSize; i++ {
			for j := 0; j < BoardSize; j++ {
				cell := b.data[i][j]
				if cell.IsSolved() || cell.MarksLength() != 1 {
					continue
				}

				solution, ok := cell.Marks.First()
				if !ok {
					return changed, errorsNewInvalidMarks(cell)
				}
				if !cell.IsValid(b, solution) {
					return changed, b.solveError()
				}

				cell.Value = solution
				cell.Marks = cell.Marks.Clear()
				changed = true
				if err := b.computeAllMarks(); err != nil {
					return changed, err
				}
				if b.isSolved() {
					return true, nil
				}
				goto nextPass
			}
		}
		return changed, nil
	nextPass:
	}
}

func errorsNewInvalidMarks(cell *Cell) error {
	return fmt.Errorf("invalid board; empty marks: Cell: %+v", cell)
}
