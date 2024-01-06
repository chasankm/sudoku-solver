package solver

import "math"

func BackTrack(data [BoardSize][BoardSize]*Cell) (bool, [BoardSize][BoardSize]*Cell) {
	solved, row, col := IsSolved(data)
	if solved {
		return true, data
	}
	for _, num := range Digits.ToArray() {
		if IsValidValue(data, row, col, Value(num)) {
			data[row][col].Value = Value(num)
			solved, _ := BackTrack(data)
			if solved {
				return true, data
			}
			data[row][col].Value = EmptyCellValue
		}
	}
	return false, [BoardSize][BoardSize]*Cell{}
}

func IsSolved(data [BoardSize][BoardSize]*Cell) (bool, int, int) {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := data[i][j]
			if cell.Value == EmptyCellValue {
				return false, i, j
			}
		}
	}
	return true, 0, 0
}

func IsValidValue(data [BoardSize][BoardSize]*Cell, row int, col int, value Value) bool {
	for _, r := range Row(data, row) {
		if r.Value == value {
			return false
		}
	}
	for _, r := range Col(data, col) {
		if r.Value == value {
			return false
		}
	}
	for _, r := range Box(data, row, col) {
		if r.Value == value {
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
		for j := 0; j < BoardSize; j++ {
			if j == index {
				cells = append(cells, data[i][j])
			}
		}
	}
	return cells
}

// Box returns the box cells in the given cell index by row and col ids
func Box(data [BoardSize][BoardSize]*Cell, rowID int, colID int) []*Cell {
	cells := make([]*Cell, 0, BoardSize)
	if (rowID < 0 || rowID >= BoardSize) || (colID < 0 || colID >= BoardSize) {
		return cells
	}
	rowMin := int(math.Floor(float64(rowID/BlockSize))) * BlockSize
	rowMax := rowMin + BlockSize
	colMin := int(math.Floor(float64(colID/BlockSize))) * BlockSize
	colMax := colMin + BlockSize
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			if i >= rowMin && i < rowMax && j >= colMin && j < colMax {
				cells = append(cells, data[i][j])
			}
			if len(cells) == BoardSize {
				return cells
			}
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
				Marks: cell.Marks.Clone(),
			}
		}
	}
	return clone
}
