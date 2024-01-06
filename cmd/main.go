package main

import (
	"fmt"
	"github.com/chasankm/sudoku-solver/pkg/solver"
	"log"
)

func main() {
	boards, parseErr := solver.ParseFile("./data/easy50.txt")
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	for index, board := range boards {
		solved, solution, seconds, solveErr := board.Solve()
		if solveErr != nil {
			givens, backTrackUsed := board.GetGivensAndBackTrack()
			fmt.Printf("[%d]: Sudoku (%d) (Backtrack: %t) could not be solved; encountered error in %.2f seconds\n", index, givens, backTrackUsed, seconds)
			fmt.Printf("Error: %s\n", solveErr.Error())
			fmt.Printf("%s\n", solution)
		} else {
			givens, backTrackUsed := board.GetGivensAndBackTrack()
			if solved {
				fmt.Printf("[%d]: Sudoku (%d) (Backtrack: %t) is solved in %.2f seconds\n", index, givens, backTrackUsed, seconds)
				fmt.Printf("%s\n", solution)
			} else {
				fmt.Printf("[%d]: Sudoku (%d) (Backtrack: %t) could not be solved in %.2f seconds\n", index, givens, backTrackUsed, seconds)
				fmt.Printf("%s\n", solution)
			}
		}
	}
}
