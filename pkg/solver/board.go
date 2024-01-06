package solver

import (
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"strconv"
	"strings"
	"time"
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

var Digits = roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9)

var ValidInputs = roaring.BitmapOf(uint32(EmptyCellValue), 1, 2, 3, 4, 5, 6, 7, 8, 9)

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

// Board is the struct of the Sudoku board
type Board struct {
	data          [BoardSize][BoardSize]*Cell
	difficulty    Difficulty
	givens        int
	backTrackUsed bool
}

// NewBoard returns new Sudoku board with the given input matrix, if there are any issues it also returns error
func NewBoard(input [BoardSize][BoardSize]Value) (*Board, error) {
	var data [BoardSize][BoardSize]*Cell
	givens := 0
	id := 0
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			value := input[i][j]
			if !ValidInputs.Contains(uint32(value)) {
				return nil, fmt.Errorf("%d is not valid input at [%d][%d]\n", value, i, j)
			}
			cell := &Cell{
				ID:    id,
				Row:   i,
				Col:   j,
				Value: input[i][j],
				Marks: roaring.NewBitmap(),
			}
			data[i][j] = cell
			id++
			if value != EmptyCellValue {
				givens++
			}
		}
	}
	if givens < MinimumGivens {
		return nil, fmt.Errorf("At least %d cells should be given to find out unique solution."+
			"Current givens: %d\n", MinimumGivens, givens)
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
	return &Board{
		data:          data,
		difficulty:    difficulty,
		givens:        givens,
		backTrackUsed: false,
	}, nil
}

// GetGivensAndBackTrack returns the givens and backTrackUsed flag
func (b *Board) GetGivensAndBackTrack() (int, bool) {
	return b.givens, b.backTrackUsed
}

// Solve is a utility function to start the solving process of given sudoku board
func (b *Board) Solve() (bool, string, float64, error) {
	begin := time.Now()

	threshold := 3
	emptyCycles := 0
	info := make(map[int]int)

	if computeErr := b.computeAllMarks(); computeErr != nil {
		return false, b.getSolution(true), time.Since(begin).Seconds(), computeErr

	}
	iteration := 0
	info[iteration] = b.emptyCells()

	for {
	solve:

		if b.hasUniqueSolutions() {
			for i := 0; i < BoardSize; i++ {
				for j := 0; j < BoardSize; j++ {
					cell := b.data[i][j]
					if !cell.IsSolved() {
						if cell.MarksLength() == 1 {
							solution := cell.Marks.ToArray()[0]
							if cell.IsValid(b, Value(solution)) {
								cell.Value = Value(solution)
								cell.Marks.Clear()
								if computeErr := b.computeAllMarks(); computeErr != nil {
									return false, b.getSolution(true), time.Since(begin).Seconds(), computeErr
								}
							}
						}
					}
				}
			}
		}

		if b.hasUniqueSolutions() {
			// There might be new solutions during the loop below, In that case, we should goto solve and continue
			goto solve
		}

		// Naked Quads strategy
		if err := b.eliminateNQ(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		if b.hasUniqueSolutions() {
			continue
		}

		// Naked Triples strategy
		if err := b.eliminateNT(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		if b.hasUniqueSolutions() {
			continue
		}

		// Naked Pairs strategy
		if err := b.eliminateNP(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		if b.hasUniqueSolutions() {
			continue
		}

		// XY Wings strategy
		s := b.totalMarks()
		if err := b.eliminateXYWings(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		diff := s - b.totalMarks()
		if b.hasUniqueSolutions() || diff > 0 {
			continue
		}

		// XYZ Wings strategy
		s = b.totalMarks()
		if err := b.eliminateXYZWings(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		diff = s - b.totalMarks()
		if b.hasUniqueSolutions() || diff > 0 {
			continue
		}

		// XWings strategy
		s = b.totalMarks()
		if err := b.eliminateXWings(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		diff = s - b.totalMarks()
		if b.hasUniqueSolutions() || diff > 0 {
			continue
		}

		// Sword Fish strategy
		s = b.totalMarks()
		if err := b.eliminateSwordFish(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		diff = s - b.totalMarks()
		if b.hasUniqueSolutions() || diff > 0 {
			continue
		}

		// Hidden Single strategy
		s = b.totalMarks()
		if err := b.eliminateHS(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		diff = s - b.totalMarks()
		if b.hasUniqueSolutions() || diff > 0 {
			continue
		}

		// Hidden Quads strategy
		s = b.totalMarks()
		if err := b.eliminateHQ(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		diff = s - b.totalMarks()
		if b.hasUniqueSolutions() {
			continue
		}

		// Hidden Triplets strategy
		s = b.totalMarks()
		if err := b.eliminateHT(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		diff = s - b.totalMarks()
		if b.hasUniqueSolutions() {
			continue
		}

		// Hidden Pairs strategy
		s = b.totalMarks()
		if err := b.eliminateHP(); err != nil {
			return false, b.getSolution(true), time.Since(begin).Seconds(), err
		}
		diff = s - b.totalMarks()
		if b.hasUniqueSolutions() || diff > 0 {
			continue
		}

		iteration++
		info[iteration] = b.emptyCells()
		if info[iteration-1] == info[iteration] {
			// Increasing the empty cycles
			emptyCycles++
		} else {
			// We have some solutions, reset the empty cycles
			emptyCycles = 0
		}
		if emptyCycles == threshold {
			// Empty cycles reached to the threshold, giving up using strategies and using backtrack
			b.backTrack()
			break
		}
		if b.isSolved() {
			break
		}
	}
	if b.isSolved() {
		return true, b.getSolution(true), time.Since(begin).Seconds(), nil
	} else {
		if b.hasInvalidMarks() {
			return false, b.getSolution(true), time.Since(begin).Seconds(), nil
		} else {
			return false, b.getSolution(true), time.Since(begin).Seconds(), nil
		}
	}
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
	return b.emptyCells() == 0
}

// solutions returns the number of marks with the given length
func (b *Board) solutions(length int) int {
	solutions := 0
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				if cell.MarksLength() == length {
					solutions++
				}
			}
		}
	}
	return solutions
}

// uniqueSolutions returns the number of unique solutions (having single candidate)
func (b *Board) uniqueSolutions() int {
	return b.solutions(1)
}

// hasUniqueSolutions returns whether there are any unique solutions or not
func (b *Board) hasUniqueSolutions() bool {
	return b.uniqueSolutions() > 0
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

// getSolution simply prints the board in well format
func (b *Board) getSolution(showInfo bool) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Difficulty: %s\n", Levels[b.difficulty]))
	builder.WriteString(fmt.Sprintf("Givens: %d\n", b.givens))
	if showInfo {
		builder.WriteString(fmt.Sprintf("BackTracking used: %t\n", b.backTrackUsed))
	}
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
	total := uint64(0)
	for _, cell := range b.unsolvedCells() {
		total += cell.Marks.GetCardinality()
	}
	return int(total)
}

// computeAllMarks simply computes all marks/candidates of each unsolved cells
func (b *Board) computeAllMarks() error {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			cell := b.data[i][j]
			if !cell.IsSolved() {
				cell.Marks = cell.ComputeCellMarks(b)
				if cell.Marks.IsEmpty() {
					fmt.Printf("%s\n", b.getSolution(true))
					return fmt.Errorf("Compute All marks: Cell: %+v\n", cell)
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
				b.data[i][j].Marks.Clear()
			}
		}
		return true
	}
	return false
}
