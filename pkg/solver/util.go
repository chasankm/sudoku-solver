package solver

import (
	"bytes"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"os"
)

// ParseFile simply parses sudoku file and returns a sudoku board for each line
func ParseFile(path string) ([]*Board, error) {
	boards := make([]*Board, 0)

	data, readFileErr := os.ReadFile(path)
	if readFileErr != nil {
		return nil, readFileErr
	}

	reader := bytes.NewReader(data)
	bufferSize := BoardSize*BoardSize + 1
	buffer := make([]byte, bufferSize)

	for {
		n, readErr := reader.Read(buffer)
		if readErr != nil {
			break
		}
		if n < bufferSize {
			continue
		}
		if !IsValidBuffer(buffer) {
			continue
		}
		var boardData [BoardSize][BoardSize]Value
		index := 0
		for i := 0; i < BoardSize; i++ {
			for j := 0; j < BoardSize; j++ {
				value := buffer[index]
				if v, ok := HexMap[value]; ok {
					boardData[i][j] = v
				}
				index++
			}
		}
		board, newBoardErr := NewBoard(boardData)
		if newBoardErr != nil {
			fmt.Printf("Board error: %s\n", newBoardErr.Error())
			continue
		}
		boards = append(boards, board)
	}
	return boards, nil
}

// FindMarksOfUnit finds marks/candidates of the given unit by getting the difference (XOR) of solved cells
// and all available digits
func FindMarksOfUnit(unit []*Cell) *roaring.Bitmap {
	solved := SolvedCells(unit)
	digits := Digits.Clone()
	digits.Xor(solved)
	return digits
}

// SolvedCells creates marks from the solved cells within the unit
func SolvedCells(cells []*Cell) *roaring.Bitmap {
	marks := roaring.NewBitmap()
	for _, cell := range cells {
		if cell.IsSolved() {
			marks.Add(uint32(cell.Value))
		}
	}
	return marks
}

// UnSolvedCells returns the unsolved cells within the unit
func UnSolvedCells(cells []*Cell) []*Cell {
	unsolved := make([]*Cell, 0)
	for _, cell := range cells {
		if !cell.IsSolved() {
			unsolved = append(unsolved, cell)
		}
	}
	return unsolved
}

// IsCellInCollection is a helper function to report whether given cell is inside of collection or not
func IsCellInCollection(cell *Cell, collection []*Cell) bool {
	for _, c := range collection {
		if cell.ID == c.ID {
			return true
		}
	}
	return false
}

// IsCellInCollections is a helper function to report whether given cell is inside of collections or not
func IsCellInCollections(cell *Cell, collection [][]*Cell) bool {
	for _, unit := range collection {
		if IsCellInCollection(cell, unit) {
			return true
		}
	}
	return false
}

// IsValidBuffer simply checks whether the current buffer holds valid board information or not
func IsValidBuffer(buffer []byte) bool {
	if len(buffer) != BoardSize*BoardSize+Offset {
		return false
	}
	if buffer[len(buffer)-1] != EOL {
		return false
	}
	for _, b := range buffer[:len(buffer)-1] {
		_, ok := HexMap[b]
		if !ok {
			return false
		}
	}
	return true
}
