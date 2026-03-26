package solver

func BackTrack(data [BoardSize][BoardSize]*Cell) (bool, [BoardSize][BoardSize]*Cell) {
	solved, valid, row, col, candidates := nextBacktrackCell(data)
	if solved {
		return true, data
	}
	if !valid {
		return false, [BoardSize][BoardSize]*Cell{}
	}
	for _, num := range candidates.ToArray() {
		value, ok := valueFromDigit(num)
		if !ok {
			return false, [BoardSize][BoardSize]*Cell{}
		}
		data[row][col].Value = value
		solved, solution := BackTrack(data)
		if solved {
			return true, solution
		}
		data[row][col].Value = EmptyCellValue
	}
	return false, [BoardSize][BoardSize]*Cell{}
}

func nextBacktrackCell(data [BoardSize][BoardSize]*Cell) (bool, bool, int, int, CandidateSet) {
	bestCount := BoardSize + 1
	bestRow, bestCol := -1, -1
	var bestMarks CandidateSet

	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := data[i][j]
			if cell.IsSolved() {
				continue
			}

			marks := candidateSetForPosition(data, i, j)
			if marks.IsEmpty() {
				return false, false, -1, -1, 0
			}

			count := marks.GetCardinality()
			if count < bestCount {
				bestCount = count
				bestRow = i
				bestCol = j
				bestMarks = marks
				if count == 1 {
					return false, true, bestRow, bestCol, bestMarks
				}
			}
		}
	}

	if bestRow == -1 {
		return !hasConflictingValues(data), !hasConflictingValues(data), -1, -1, 0
	}

	return false, true, bestRow, bestCol, bestMarks
}

func IsValidValue(data [BoardSize][BoardSize]*Cell, row int, col int, value Value) bool {
	for _, peerID := range peerIDs[cellID(row, col)] {
		if cellByID(data, peerID).Value == value {
			return false
		}
	}
	return true
}

// Row returns the row in the given index
func Row(data [BoardSize][BoardSize]*Cell, index int) []*Cell {
	cells := make([]*Cell, 0, BoardSize)
	if index < 0 || index >= BoardSize {
		return cells
	}
	return data[index][:]
}

// Col returns the col in the given index
func Col(data [BoardSize][BoardSize]*Cell, index int) []*Cell {
	cells := make([]*Cell, 0, BoardSize)
	if index < 0 || index >= BoardSize {
		return cells
	}
	for i := 0; i < BoardSize; i++ {
		cells = append(cells, data[i][index])
	}
	return cells
}

// Box returns the box cells in the given cell index by row and col ids
func Box(data [BoardSize][BoardSize]*Cell, rowID int, colID int) []*Cell {
	cells := make([]*Cell, 0, BoardSize)
	if (rowID < 0 || rowID >= BoardSize) || (colID < 0 || colID >= BoardSize) {
		return cells
	}
	rowMin := (rowID / BlockSize) * BlockSize
	colMin := (colID / BlockSize) * BlockSize
	for i := rowMin; i < rowMin+BlockSize; i++ {
		for j := colMin; j < colMin+BlockSize; j++ {
			cells = append(cells, data[i][j])
		}
	}
	return cells
}

func CloneData(data [BoardSize][BoardSize]*Cell) [BoardSize][BoardSize]*Cell {
	var clone [BoardSize][BoardSize]*Cell
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := data[i][j]
			clone[i][j] = &Cell{
				ID:    cell.ID,
				Row:   cell.Row,
				Col:   cell.Col,
				Value: cell.Value,
				Marks: cell.Marks,
			}
		}
	}
	return clone
}
