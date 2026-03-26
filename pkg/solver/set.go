package solver

import (
	"math/bits"
	"strconv"
	"strings"
)

type CandidateSet uint16

func CandidateSetOf(values ...int) CandidateSet {
	var set CandidateSet
	for _, value := range values {
		if value >= 1 && value <= BoardSize {
			set |= 1 << value
		}
	}
	return set
}

func (s CandidateSet) Clone() CandidateSet {
	return s
}

func (s CandidateSet) Contains(value int) bool {
	if value == 0 || value > BoardSize {
		return false
	}
	return s&(1<<value) != 0
}

func (s CandidateSet) Add(value int) CandidateSet {
	if value >= 1 && value <= BoardSize {
		s |= 1 << value
	}
	return s
}

func (s CandidateSet) Clear() CandidateSet {
	return 0
}

func (s CandidateSet) IsEmpty() bool {
	return s == 0
}

func (s CandidateSet) GetCardinality() int {
	return bits.OnesCount16(uint16(s))
}

func (s CandidateSet) ToArray() []int {
	values := make([]int, 0, s.GetCardinality())
	for value := 1; value <= BoardSize; value++ {
		if s.Contains(value) {
			values = append(values, value)
		}
	}
	return values
}

func (s CandidateSet) First() (Value, bool) {
	for value := 1; value <= BoardSize; value++ {
		if s.Contains(value) {
			return valueFromDigit(value)
		}
	}
	return EmptyCellValue, false
}

func (s CandidateSet) And(other CandidateSet) CandidateSet {
	return s & other
}

func (s CandidateSet) AndNot(other CandidateSet) CandidateSet {
	return s &^ other
}

func (s CandidateSet) Xor(other CandidateSet) CandidateSet {
	return s ^ other
}

func (s CandidateSet) String() string {
	values := s.ToArray()
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, strconv.Itoa(value))
	}
	return "{" + strings.Join(parts, ",") + "}"
}

// ParIntersect returns the intersection of all given units.
func ParIntersect(units ...CandidateSet) CandidateSet {
	if len(units) == 0 {
		return 0
	}
	result := units[0]
	for _, unit := range units[1:] {
		result &= unit
	}
	return result
}

// ParIntersectCells returns the intersection of all given cells.
func ParIntersectCells(cells []*Cell) CandidateSet {
	if len(cells) == 0 {
		return 0
	}
	result := cells[0].Marks
	for _, cell := range cells[1:] {
		result &= cell.Marks
	}
	return result
}

// ParUnion returns the union of all given units.
func ParUnion(units ...CandidateSet) CandidateSet {
	var result CandidateSet
	for _, unit := range units {
		result |= unit
	}
	return result
}

// ParUnionCells returns the union of all given cells.
func ParUnionCells(cells []*Cell) CandidateSet {
	var result CandidateSet
	for _, cell := range cells {
		result |= cell.Marks
	}
	return result
}
