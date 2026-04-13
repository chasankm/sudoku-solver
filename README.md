# sudoku-solver

`sudoku-solver` is a Go Sudoku solver that combines human-style candidate elimination strategies with backtracking when logic alone stalls.

It can be used as:

- A CLI that parses puzzle files and solves boards in parallel
- A Go package via `github.com/chasankm/sudoku-solver/pkg/solver`

The project is now pure Go and does not depend on external bitmap libraries.

## Requirements

- Go 1.25+

## What It Does

For each board, the solver:

1. Validates the initial givens
2. Initializes candidates for unsolved cells
3. Resolves singles
4. Applies advanced strategies in a fixed order
5. Falls back to backtracking if the logical passes no longer make progress

Backtracking uses a minimum-remaining-values style choice by selecting the unsolved cell with the fewest candidates first.

## Implemented Strategies

The current strategy pipeline is:

- Naked Quads
- Naked Triples
- Naked Pairs
- Hidden Single
- Locked Candidates
- XY Wings
- XYZ Wings
- X Wings
- Sword Fish
- Hidden Quads
- Hidden Triplets
- Hidden Pairs

The solve response records which strategies were used for each puzzle and whether backtracking was required.

## Input Format

`ParseFile` reads one board per line.

Supported cell values:

- `1`-`9` for givens
- `0` or `.` for empty cells

Rules:

- Each line must contain exactly 81 board characters
- Invalid lines are skipped
- Boards with conflicting givens are skipped
- Boards with fewer than 17 givens are rejected

Sample datasets are included in [data/easy50.txt](https://github.com/chasankm/sudoku-solver/blob/main/data/easy50.txt) and [data/top95.txt](https://github.com/chasankm/sudoku-solver/blob/main/data/top95.txt).

## Run The CLI

From the repository root:

```bash
go run ./cmd
```

The current CLI reads `./data/top95.txt`, solves boards in batches sized to `runtime.NumCPU()`, and prints:

- Board index
- Estimated difficulty from clue count
- Number of givens
- Solved status
- Whether backtracking was used
- Strategies used
- Solve duration
- Initial board state
- Final board state

Example:

```text
Number of workers: 16
Number of boards parsed 95
16 of 95 boards have been processed

Index: 0
Difficulty: Evil
Givens: 17
Is Solved: true
BackTracking used: false
Strategies used: [Naked Triples, Naked Pairs, Naked Quads, Hidden Single, Locked Candidates]
Duration: 0.32 seconds
Initial state:
*_______*_______*______*
| 4 _ _ | _ _ _ | 8 _ 5
| _ 3 _ | _ _ _ | _ _ _
| _ _ _ | 7 _ _ | _ _ _
*_______*_______*______*
| _ 2 _ | _ _ _ | _ 6 _
| _ _ _ | _ 8 _ | 4 _ _
| _ _ _ | _ 1 _ | _ _ _
*_______*_______*______*
| _ _ _ | 6 _ 3 | _ 7 _
| 5 _ _ | 2 _ _ | _ _ _
| 1 _ 4 | _ _ _ | _ _ _

Solution:
*_______*_______*______*
| 4 1 7 | 3 6 9 | 8 2 5
| 6 3 2 | 1 5 8 | 9 4 7
| 9 5 8 | 7 2 4 | 3 1 6
...
```

## Use As A Library

```go
package main

import (
	"fmt"
	"log"

	"github.com/chasankm/sudoku-solver/pkg/solver"
)

func main() {
	boards, err := solver.ParseFile("./data/easy50.txt")
	if err != nil {
		log.Fatal(err)
	}

	for i, board := range boards {
		result := board.Solve()
		fmt.Printf("Board %d solved=%t backtracking=%t\n", i, result.IsSolved, result.BackTrackingUsed)
		if result.Error != nil {
			fmt.Println(result.Error)
		}
		fmt.Println(result.Print())
	}
}
```

`Solve()` returns a `SolveResponse` with:

- `Difficulty`
- `Givens`
- `InitialState`
- `Solution`
- `Duration`
- `IsSolved`
- `BackTrackingUsed`
- `StrategiesUsed`
- `Error`

## Validation Behavior

The solver rejects invalid starting states early and also reports failures when a board reaches an inconsistent or unsolved terminal state. Typical failure reasons are:

- Conflicting givens
- Invalid candidate state
- No solution found

## Tests

Run the test suite with:

```bash
go test ./...
```

The repository also includes a benchmark for the bundled `top95` dataset:

```bash
go test -bench=. ./pkg/solver
```
