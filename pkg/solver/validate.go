package solver

import "errors"

func hasConflictingValues(data [BoardSize][BoardSize]*Cell) bool {
	for i := 0; i < BoardSize; i++ {
		if unitHasConflictingValues(Row(data, i)) || unitHasConflictingValues(Col(data, i)) {
			return true
		}
	}
	for row := 0; row < BoardSize; row += BlockSize {
		for col := 0; col < BoardSize; col += BlockSize {
			if unitHasConflictingValues(Box(data, row, col)) {
				return true
			}
		}
	}
	return false
}

func unitHasConflictingValues(unit []*Cell) bool {
	seen := make(map[Value]struct{}, BoardSize)
	for _, cell := range unit {
		if !cell.IsSolved() {
			continue
		}
		if _, ok := seen[cell.Value]; ok {
			return true
		}
		seen[cell.Value] = struct{}{}
	}
	return false
}

func (b *Board) validateForSolve() error {
	if !b.isValid() {
		return errors.New("invalid board; conflicting values found")
	}
	return nil
}

func (b *Board) isValid() bool {
	return !hasConflictingValues(b.data)
}

func (b *Board) solveError() error {
	if b.hasInvalidMarks() {
		return errors.New("invalid board; some cells have invalid marks")
	}
	if !b.isValid() {
		return errors.New("invalid board; conflicting values found")
	}
	return errors.New("unable to find a solution")
}
