package solver

import "time"

const stalledCycleThreshold = 1

// Solve is a utility function to start the solving process of given sudoku board.
func (b *Board) Solve() *SolveResponse {
	begin := time.Now()

	if err := b.validateForSolve(); err != nil {
		return b.buildSolveResponse(begin, err)
	}
	if err := b.initializeCandidates(); err != nil {
		return b.buildSolveResponse(begin, err)
	}
	if b.isSolved() {
		return b.buildSolveResponse(begin, nil)
	}

	stalledCycles := 0
	for !b.isSolved() {
		changed, err := b.resolveSingles()
		if err != nil {
			return b.buildSolveResponse(begin, err)
		}
		if b.isSolved() {
			break
		}
		if changed {
			stalledCycles = 0
			continue
		}

		changed, err = b.applyStrategies()
		if err != nil {
			return b.buildSolveResponse(begin, err)
		}
		if b.isSolved() {
			break
		}
		if changed {
			stalledCycles = 0
			continue
		}

		stalledCycles++
		if stalledCycles >= stalledCycleThreshold {
			b.backTrack()
			break
		}
	}

	return b.buildSolveResponse(begin, nil)
}

func (b *Board) buildSolveResponse(begin time.Time, err error) *SolveResponse {
	if err == nil && !b.isSolved() {
		err = b.solveError()
	}

	return &SolveResponse{
		Difficulty:       Levels[b.difficulty],
		Givens:           b.givens,
		InitialState:     b.initialState,
		Solution:         b.getState(),
		Duration:         time.Since(begin).Seconds(),
		IsSolved:         err == nil,
		BackTrackingUsed: b.backTrackUsed,
		StrategiesUsed:   b.strategiesUsed,
		Error:            err,
	}
}
