package solver

import (
	"fmt"
	"strings"
)

type SolveResponse struct {
	Difficulty       string
	Givens           int
	InitialState     string
	Solution         string
	Duration         float64
	IsSolved         bool
	BackTrackingUsed bool
	StrategiesUsed   []string
	Error            error
}

func (r *SolveResponse) Print() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Difficulty: %s\n", r.Difficulty))
	builder.WriteString(fmt.Sprintf("Givens: %d\n", r.Givens))
	builder.WriteString(fmt.Sprintf("Is Solved: %t\n", r.IsSolved))
	builder.WriteString(fmt.Sprintf("BackTracking used: %t\n", r.BackTrackingUsed))
	builder.WriteString(fmt.Sprintf("Strategies used: [%s]\n", strings.Join(r.StrategiesUsed, ", ")))
	builder.WriteString(fmt.Sprintf("Duration: %.2f seconds\n", r.Duration))
	if r.Error != nil {
		builder.WriteString(fmt.Sprintf("Error: %s\n", r.Error.Error()))
	}
	builder.WriteString(fmt.Sprintf("Initial state: \n%s\n", r.InitialState))
	if r.IsSolved {
		builder.WriteString(fmt.Sprintf("Solution: \n%s\n", r.Solution))
	} else {
		builder.WriteString(fmt.Sprintf("Last state: \n%s\n", r.Solution))
	}

	return builder.String()
}
