package solver

import (
	"os"
	"path/filepath"
	"testing"
)

const solvedBoard = "483921657967345821251876493548132976729564138136798245372689514814253769695417382"

func TestNewBoardRejectsConflictingGivens(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "row conflict",
			input: "113020600900305001001806400008102900700000008006708200002609500800203009005010300",
		},
		{
			name:  "column conflict",
			input: "003020600900305001001806400008102900700000008006708200002609500800203009905010300",
		},
		{
			name:  "box conflict",
			input: "103020600900305001001806400008102900700000008006708200002609500800203009005010300",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBoard(mustGridFromString(t, tt.input))
			if err == nil {
				t.Fatal("expected conflicting givens to be rejected")
			}
		})
	}
}

func TestSolveSolvesKnownEasyPuzzle(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, "003020600900305001001806400008102900700000008006708200002609500800203009005010300"))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	response := board.Solve()
	if response.Error != nil {
		t.Fatalf("Solve() error = %v", response.Error)
	}
	if !response.IsSolved {
		t.Fatal("Solve() did not solve the puzzle")
	}

	expected, err := NewBoard(mustGridFromString(t, "483921657967345821251876493548132976729564138136798245372689514814253769695417382"))
	if err != nil {
		t.Fatalf("NewBoard(expected) error = %v", err)
	}

	if response.Solution != expected.getState() {
		t.Fatalf("Solve() produced unexpected solution:\n%s", response.Solution)
	}
}

func TestSolveRejectsInvalidCompletedState(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, solvedBoard))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	board.data[0][1].Value = board.data[0][0].Value

	response := board.Solve()
	if response.Error == nil {
		t.Fatal("Solve() expected an error for an invalid completed board")
	}
	if response.IsSolved {
		t.Fatal("Solve() reported an invalid completed board as solved")
	}
}

func TestParseFileSkipsInvalidLinesAndParsesValidBoards(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "boards.txt")
	content := "" +
		"003020600900305001001806400008102900700000008006708200002609500800203009005010300\n" +
		"not-a-valid-board-line\n" +
		"200080300060070084030500209000105408000000000402706000301007040720040060004010003\n"

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	boards, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	if len(boards) != 2 {
		t.Fatalf("ParseFile() parsed %d boards, want 2", len(boards))
	}

	givens, backtrackUsed := boards[0].GetGivensAndBackTrack()
	if givens != 32 {
		t.Fatalf("first board givens = %d, want 32", givens)
	}
	if backtrackUsed {
		t.Fatal("newly parsed board should not report backtracking")
	}
}

func TestNextBacktrackCellChoosesMinimumCandidateCell(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, "003020600900305001001806400008102900700000008006708200002609500800203009005010300"))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	solved, valid, row, col, marks := nextBacktrackCell(board.data)
	if solved {
		t.Fatal("expected puzzle to be unsolved")
	}
	if !valid {
		t.Fatal("expected puzzle state to be valid")
	}
	if marks.IsEmpty() {
		t.Fatal("expected MRV cell to have candidates")
	}

	best := marks.GetCardinality()
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := board.data[i][j]
			if cell.IsSolved() {
				continue
			}
			candidates := candidateSetForPosition(board.data, i, j)
			if candidates.GetCardinality() < best {
				t.Fatalf("selected cell [%d][%d] had %d candidates, but [%d][%d] had %d", row, col, best, i, j, candidates.GetCardinality())
			}
		}
	}
}

func TestOrderedStrategiesUseExpectedOrder(t *testing.T) {
	expected := []StrategyName{
		NakedQuadsStrategy,
		NakedTriplesStrategy,
		NakedPairsStrategy,
		HiddenSingleStrategy,
		LockedCandidatesStrategy,
		XYWingsStrategy,
		XYZWingsStrategy,
		XWingsStrategy,
		SwordFishStrategy,
		HiddenQuadsStrategy,
		HiddenTripletsStrategy,
		HiddenPairsStrategy,
	}

	if len(orderedStrategies) != len(expected) {
		t.Fatalf("orderedStrategies length = %d, want %d", len(orderedStrategies), len(expected))
	}

	for i, strategy := range orderedStrategies {
		if strategy.Name() != expected[i] {
			t.Fatalf("orderedStrategies[%d] = %q, want %q", i, strategy.Name(), expected[i])
		}
	}
}

func TestEliminateNakedPairsFixture(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, solvedBoard))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	setCandidates(board, 0, 0, 1, 2)
	setCandidates(board, 0, 1, 1, 2)
	setCandidates(board, 0, 2, 1, 2, 3)

	if err := EliminateNakedPairs([][]*Cell{board.row(0)}); err != nil {
		t.Fatalf("EliminateNakedPairs() error = %v", err)
	}

	if board.data[0][2].Marks != CandidateSetOf(3) {
		t.Fatalf("target marks = %s, want {3}", board.data[0][2].Marks.String())
	}
}

func TestEliminateHiddenPairsFixture(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, solvedBoard))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	setCandidates(board, 0, 0, 1, 2, 3)
	setCandidates(board, 0, 1, 1, 2, 4)
	setCandidates(board, 0, 2, 3, 4, 5)
	setCandidates(board, 0, 3, 3, 5)

	if err := EliminateHiddenPairs([][]*Cell{board.row(0)}); err != nil {
		t.Fatalf("EliminateHiddenPairs() error = %v", err)
	}

	if board.data[0][0].Marks != CandidateSetOf(1, 2) {
		t.Fatalf("first pair marks = %s, want {1,2}", board.data[0][0].Marks.String())
	}
	if board.data[0][1].Marks != CandidateSetOf(1, 2) {
		t.Fatalf("second pair marks = %s, want {1,2}", board.data[0][1].Marks.String())
	}
}

func TestEliminateLockedCandidatesPointingPair(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, solvedBoard))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	setCandidates(board, 0, 0, 5, 1)
	setCandidates(board, 0, 1, 5, 2)
	setCandidates(board, 0, 3, 5, 6)

	if err := board.eliminateLockedCandidates(); err != nil {
		t.Fatalf("eliminateLockedCandidates() error = %v", err)
	}

	if board.data[0][3].Marks.Contains(5) {
		t.Fatal("pointing pair did not eliminate candidate 5 from the row outside the box")
	}
}

func TestEliminateLockedCandidatesClaimingPair(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, solvedBoard))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	setCandidates(board, 0, 0, 7, 1)
	setCandidates(board, 0, 1, 7, 2)
	setCandidates(board, 1, 2, 7, 3)

	if err := board.eliminateLockedCandidates(); err != nil {
		t.Fatalf("eliminateLockedCandidates() error = %v", err)
	}

	if board.data[1][2].Marks.Contains(7) {
		t.Fatal("claiming pair did not eliminate candidate 7 from the box outside the row")
	}
}

func TestEliminateXWingsFindsColumnBasedPattern(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, solvedBoard))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	setCandidates(board, 0, 0, 5, 6)
	setCandidates(board, 3, 0, 5, 6)
	setCandidates(board, 0, 3, 5, 7)
	setCandidates(board, 3, 3, 5, 7)
	setCandidates(board, 0, 7, 5, 8)

	if err := board.eliminateXWings(); err != nil {
		t.Fatalf("eliminateXWings() error = %v", err)
	}

	if board.data[0][7].Marks.Contains(5) {
		t.Fatal("column-based X-Wing did not eliminate candidate 5 from the target row")
	}
}

func TestEliminateSwordFishFindsColumnBasedPattern(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, solvedBoard))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	setCandidates(board, 0, 0, 5, 6)
	setCandidates(board, 3, 0, 5, 6)
	setCandidates(board, 0, 3, 5, 7)
	setCandidates(board, 6, 3, 5, 7)
	setCandidates(board, 3, 6, 5, 8)
	setCandidates(board, 6, 6, 5, 8)
	setCandidates(board, 0, 7, 5, 9)
	setCandidates(board, 0, 8, 5, 2)
	setCandidates(board, 3, 8, 5, 2)

	if err := board.eliminateSwordFish(); err != nil {
		t.Fatalf("eliminateSwordFish() error = %v", err)
	}

	if board.data[0][7].Marks.Contains(5) {
		t.Fatal("column-based SwordFish did not eliminate candidate 5 from the target row")
	}
}

func TestSolveTop95Board61DoesNotFailInXYWing(t *testing.T) {
	board, err := NewBoard(mustGridFromString(t, "1.....3.8.6.4..............2.3.1...........758.........7.5...6.....8.2...4......."))
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	response := board.Solve()
	if response.Error != nil {
		t.Fatalf("Solve() error = %v", response.Error)
	}
	if !response.IsSolved {
		t.Fatal("Solve() did not solve top95 board 61")
	}
}

func mustGridFromString(t *testing.T, input string) [BoardSize][BoardSize]Value {
	t.Helper()

	if len(input) != BoardSize*BoardSize {
		t.Fatalf("invalid board length %d", len(input))
	}

	var grid [BoardSize][BoardSize]Value
	for i := 0; i < len(input); i++ {
		value, ok := HexMap[input[i]]
		if !ok {
			t.Fatalf("invalid board character %q at index %d", input[i], i)
		}
		grid[i/BoardSize][i%BoardSize] = value
	}
	return grid
}

func setCandidates(board *Board, row int, col int, marks ...int) {
	board.data[row][col].Value = EmptyCellValue
	board.data[row][col].Marks = CandidateSetOf(marks...)
}
