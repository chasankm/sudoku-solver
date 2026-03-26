package solver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	BoardSize      = 9
	BlockSize      = 3
	EmptyCellValue = Value(0)
	MinimumGivens  = 17
	Offset         = 1
	EOL            = 0x0A
)

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
	Expert
	Evil
)

var Levels = map[Difficulty]string{
	Easy:   "Easy",
	Medium: "Medium",
	Hard:   "Hard",
	Expert: "Expert",
	Evil:   "Evil",
}

type StrategyName string

const (
	NakedQuadsStrategy       StrategyName = "Naked Quads"
	NakedTriplesStrategy     StrategyName = "Naked Triples"
	NakedPairsStrategy       StrategyName = "Naked Pairs"
	LockedCandidatesStrategy StrategyName = "Locked Candidates"
	XYWingsStrategy          StrategyName = "XY Wings"
	XYZWingsStrategy         StrategyName = "XYZ Wings"
	XWingsStrategy           StrategyName = "X Wings"
	SwordFishStrategy        StrategyName = "Sword Fish"
	HiddenSingleStrategy     StrategyName = "Hidden Single"
	HiddenQuadsStrategy      StrategyName = "Hidden Quads"
	HiddenTripletsStrategy   StrategyName = "Hidden Triplets"
	HiddenPairsStrategy      StrategyName = "Hidden Pairs"
)

func (s StrategyName) String() string {
	return string(s)
}

var Digits = CandidateSetOf(1, 2, 3, 4, 5, 6, 7, 8, 9)

var HexMap = map[byte]Value{
	0x2e: 0x00,
	0x30: 0x00,
	0x31: 0x01,
	0x32: 0x02,
	0x33: 0x03,
	0x34: 0x04,
	0x35: 0x05,
	0x36: 0x06,
	0x37: 0x07,
	0x38: 0x08,
	0x39: 0x09,
}

var digitValues = [...]Value{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

func valueFromDigit(digit int) (Value, bool) {
	if digit < 1 || digit > BoardSize {
		return EmptyCellValue, false
	}
	return digitValues[digit], true
}

// Board is the struct of the Sudoku board
type Board struct {
	data           [BoardSize][BoardSize]*Cell
	initialState   string
	difficulty     Difficulty
	givens         int
	backTrackUsed  bool
	strategiesUsed []string
}

// NewBoard returns new Sudoku board with the given input matrix, if there are any issues it also returns error
func NewBoard(input [BoardSize][BoardSize]Value) (*Board, error) {
	var data [BoardSize][BoardSize]*Cell
	givens := 0
	id := 0
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			value := input[i][j]
			if value > Value(BoardSize) {
				return nil, fmt.Errorf("%d is not valid input at [%d][%d]", value, i, j)
			}
			cell := &Cell{
				ID:    id,
				Row:   i,
				Col:   j,
				Value: input[i][j],
				Marks: 0,
			}
			data[i][j] = cell
			id++
			if value != EmptyCellValue {
				givens++
			}
		}
	}
	if hasConflictingValues(data) {
		return nil, errors.New("board contains conflicting givens")
	}
	if givens < MinimumGivens {
		return nil, fmt.Errorf("at least %d cells should be given to find out unique solution; current givens: %d", MinimumGivens, givens)
	}
	var difficulty Difficulty
	if givens > 32 {
		difficulty = Easy
	}
	if givens >= 30 && givens <= 32 {
		difficulty = Medium
	}
	if givens >= 28 && givens < 30 {
		difficulty = Hard
	}
	if givens >= 23 && givens < 28 {
		difficulty = Expert
	}
	if givens < 23 {
		difficulty = Evil
	}
	board := &Board{
		data:           data,
		initialState:   "",
		difficulty:     difficulty,
		givens:         givens,
		backTrackUsed:  false,
		strategiesUsed: make([]string, 0),
	}
	// Storing the initial state before Solve method is called
	board.initialState = board.getState()

	return board, nil
}

// GetGivensAndBackTrack returns the givens and backTrackUsed flag
func (b *Board) GetGivensAndBackTrack() (int, bool) {
	return b.givens, b.backTrackUsed
}

// addStrategy
func (b *Board) addStrategy(strategy StrategyName) {
	for _, s := range b.strategiesUsed {
		if s == strategy.String() {
			// This strategy is already added, returning
			return
		}
	}
	// New strategy adding it
	b.strategiesUsed = append(b.strategiesUsed, strategy.String())
}

// row returns the row in the given index
func (b *Board) row(index int) []*Cell {
	return Row(b.data, index)
}

// col returns the col in the given index
func (b *Board) col(index int) []*Cell {
	return Col(b.data, index)
}

// box returns the box cells in the given cell index by row and col ids
func (b *Board) box(rowID int, colID int) []*Cell {
	return Box(b.data, rowID, colID)
}

// emptyCells returns the number of the unsolved cells
func (b *Board) emptyCells() int {
	empty := 0
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				empty++
			}
		}
	}
	return empty
}

// isSolved simply returns whether the board is solved or not
func (b *Board) isSolved() bool {
	return b.emptyCells() == 0 && b.isValid()
}

// hasInvalidMarks returns whether board have invalid marks or not. This happens if the initial board is wrong
func (b *Board) hasInvalidMarks() bool {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			// Cell is not solved and having zero marks
			if !cell.IsSolved() && cell.MarksLength() == 0 {
				return true
			}
		}
	}
	return false
}

// getState returns the current state of the board
func (b *Board) getState() string {
	var builder strings.Builder
	for i := 0; i < BoardSize; i++ {
		if i%BlockSize == 0 {
			builder.WriteString("*_______*_______*______*\n")
		}
		row := b.row(i)
		for j := 0; j < len(row); j++ {
			cell := row[j]
			if j%BlockSize == 0 {
				builder.WriteString("| ")
			}
			if !cell.IsSolved() {
				builder.WriteString("_ ")
			} else {
				builder.WriteString(strconv.Itoa(int(cell.Value)) + " ")
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// unsolvedCells simply returns all unsolved cells within a slice
func (b *Board) unsolvedCells() []*Cell {
	unsolved := make([]*Cell, 0)
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				unsolved = append(unsolved, cell)
			}
		}
	}
	return unsolved
}

// totalMarks returns the total marks/candidates of the unsolved cells
func (b *Board) totalMarks() int {
	total := 0
	for _, cell := range b.unsolvedCells() {
		total += cell.Marks.GetCardinality()
	}
	return total
}

// computeAllMarks simply computes all marks/candidates of each unsolved cells
func (b *Board) computeAllMarks() error {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				cell.Marks = cell.ComputeCellMarks(b)
				if cell.Marks.IsEmpty() {
					return fmt.Errorf("compute all marks: cell: %+v", cell)
				}
			}
		}
	}
	return nil
}

// eliminateNP simply eliminates marks/candidates using naked pair elimination strategy for each unsolved cells
func (b *Board) eliminateNP() error {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				if eliminateErr := EliminateNakedPairs(cell.CellUnits(b)); eliminateErr != nil {
					return eliminateErr
				}
			}
		}
	}
	return nil
}

// eliminateNT simply eliminates marks/candidates using naked triple elimination strategy for each unsolved cells
func (b *Board) eliminateNT() error {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				if eliminateErr := EliminateNakedTriplets(cell.CellUnits(b)); eliminateErr != nil {
					return eliminateErr
				}
			}
		}
	}
	return nil
}

// eliminateNQ simply eliminates marks/candidates using naked quad elimination strategy for each unsolved cells
func (b *Board) eliminateNQ() error {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				if eliminateErr := EliminateNakedQuads(cell.CellUnits(b)); eliminateErr != nil {
					return eliminateErr
				}
			}
		}
	}
	return nil
}

// eliminateHS simply eliminates marks/candidates using hidden single elimination strategy for each unsolved cells
func (b *Board) eliminateHS() error {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				if eliminateErr := EliminateHiddenSingles(cell.CellUnits(b)); eliminateErr != nil {
					return eliminateErr
				}
			}
		}
	}
	return nil
}

// eliminateHP simply eliminates marks/candidates using hidden pair eliminations strategy for each unsolved cells
func (b *Board) eliminateHP() error {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				if eliminateErr := EliminateHiddenPairs(cell.CellUnits(b)); eliminateErr != nil {
					return eliminateErr
				}
			}
		}
	}
	return nil
}

// eliminateHT simply eliminates marks/candidates using hidden triplet eliminations strategy for each unsolved cells
func (b *Board) eliminateHT() error {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				if eliminateErr := EliminateHiddenTriplets(cell.CellUnits(b)); eliminateErr != nil {
					return eliminateErr
				}
			}
		}
	}
	return nil
}

// eliminateHQ simply eliminates marks/candidates using hidden quads eliminations strategy for each unsolved cells
func (b *Board) eliminateHQ() error {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				if eliminateErr := EliminateHiddenQuads(cell.CellUnits(b)); eliminateErr != nil {
					return eliminateErr
				}
			}
		}
	}
	return nil
}

// eliminateLockedCandidates simply eliminates marks/candidates using pointing/claiming intersections.
func (b *Board) eliminateLockedCandidates() error {
	return EliminateLockedCandidates(b)
}

// eliminateXYWings simply eliminates marks/candidates using XY Wings strategy for the board
func (b *Board) eliminateXYWings() error {
	return EliminateXYWings(b.unsolvedCells(), b)
}

// eliminateXYZWings simply eliminates marks/candidates using XYZ Wings strategy for the board
func (b *Board) eliminateXYZWings() error {
	return EliminateXYZWings(b.unsolvedCells(), b)
}

// eliminateXWings simply eliminates marks/candidates using X Wings strategy for the board
func (b *Board) eliminateXWings() error {
	return EliminateXWings(b)
}

// eliminateSwordFish simply eliminates marks/candidates using Sword Fish strategy for the board
func (b *Board) eliminateSwordFish() error {
	return EliminateSwordFish(b)
}

// backTrack simply tries to find out a unique solution where strategies no more producing solutions or eliminating candidates
func (b *Board) backTrack() bool {
	clone := CloneData(b.data)
	solved, solution := BackTrack(clone)
	if solved {
		b.backTrackUsed = true
		for i := 0; i < BoardSize; i++ {
			for j := 0; j < BoardSize; j++ {
				b.data[i][j].Value = solution[i][j].Value
				b.data[i][j].Marks = b.data[i][j].Marks.Clear()
			}
		}
		return true
	}
	return false
}
