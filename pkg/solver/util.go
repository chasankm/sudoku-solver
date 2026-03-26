package solver

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// ParseFile simply parses sudoku file and returns a sudoku board for each line
func ParseFile(path string) ([]*Board, error) {
	boards := make([]*Board, 0)

	cleanPath := filepath.Clean(path)
	root, err := os.OpenRoot(filepath.Dir(cleanPath))
	if err != nil {
		return nil, err
	}
	defer func(root *os.Root) {
		if cErr := root.Close(); cErr != nil {
			fmt.Printf("Error closing root: %s\n", cErr.Error())
		}
	}(root)

	file, err := root.Open(filepath.Base(cleanPath))
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if cErr := file.Close(); cErr != nil {
			fmt.Printf("Error closing file: %s\n", cErr.Error())
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		buffer := append(append(make([]byte, 0, BoardSize*BoardSize+Offset), scanner.Bytes()...), EOL)
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
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return boards, nil
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
