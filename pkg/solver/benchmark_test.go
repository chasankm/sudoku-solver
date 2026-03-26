package solver

import "testing"

var benchmarkSolvedBoards int

func BenchmarkSolveTop95(b *testing.B) {
	for i := 0; i < b.N; i++ {
		boards, err := ParseFile("../../data/top95.txt")
		if err != nil {
			b.Fatalf("ParseFile() error = %v", err)
		}
		solved := 0
		for _, board := range boards {
			response := board.Solve()
			if response.IsSolved {
				solved++
			}
		}
		benchmarkSolvedBoards = solved
	}
}
