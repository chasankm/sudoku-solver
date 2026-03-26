package solver

import (
	"strconv"
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
	builder.WriteString("Difficulty: ")
	builder.WriteString(r.Difficulty)
	builder.WriteByte('\n')
	builder.WriteString("Givens: ")
	builder.WriteString(strconv.Itoa(r.Givens))
	builder.WriteByte('\n')
	builder.WriteString("Is Solved: ")
	builder.WriteString(strconv.FormatBool(r.IsSolved))
	builder.WriteByte('\n')
	builder.WriteString("BackTracking used: ")
	builder.WriteString(strconv.FormatBool(r.BackTrackingUsed))
	builder.WriteByte('\n')
	builder.WriteString("Strategies used: [")
	builder.WriteString(strings.Join(r.StrategiesUsed, ", "))
	builder.WriteString("]\n")
	builder.WriteString("Duration: ")
	builder.WriteString(strconv.FormatFloat(r.Duration, 'f', 2, 64))
	builder.WriteString(" seconds\n")
	if r.Error != nil {
		builder.WriteString("Error: ")
		builder.WriteString(r.Error.Error())
		builder.WriteByte('\n')
	}
	builder.WriteString("Initial state: \n")
	builder.WriteString(r.InitialState)
	builder.WriteByte('\n')
	if r.IsSolved {
		builder.WriteString("Solution: \n")
	} else {
		builder.WriteString("Last state: \n")
	}
	builder.WriteString(r.Solution)
	builder.WriteByte('\n')

	return builder.String()
}
