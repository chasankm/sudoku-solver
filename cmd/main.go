package main

import (
	"fmt"
	"github.com/chasankm/sudoku-solver/pkg/solver"
	"log"
	"runtime"
	"slices"
	"sync"
	"time"
)

type Results struct {
	Solutions map[int]*solver.SolveResponse
	Mutex     sync.Mutex
}

func (r *Results) AddSolution(index int, solution *solver.SolveResponse) {
	r.Mutex.Lock()
	r.Solutions[index] = solution
	r.Mutex.Unlock()
}

func (r *Results) PrintAll() {
	indexes := make([]int, 0, len(r.Solutions))
	for index := range r.Solutions {
		indexes = append(indexes, index)
	}
	slices.Sort(indexes)
	for _, index := range indexes {
		solution, ok := r.Solutions[index]
		if ok {
			fmt.Printf("Index: %d\n", index)
			fmt.Printf("%s", solution.Print())
		}
	}
}

func (r *Results) NumberOfUnsolved() int {
	total := 0
	for _, v := range r.Solutions {
		if v.IsSolved == false {
			total += 1
		}
	}
	return total
}

func (r *Results) PrintUnsolved() {
	indexes := make([]int, 0, len(r.Solutions))
	for index, solution := range r.Solutions {
		if solution.IsSolved == false {
			indexes = append(indexes, index)
		}
	}
	slices.Sort(indexes)
	for _, index := range indexes {
		solution, ok := r.Solutions[index]
		if ok {
			fmt.Printf("Index: %d\n", index)
			fmt.Printf("%s", solution.Print())
		}
	}
}

func SolveBoardsInBatch(results *Results, boards map[int]*solver.Board) int {
	var wg sync.WaitGroup

	for idx, b := range boards {
		wg.Add(1)
		go func(index int, results *Results, board *solver.Board, w *sync.WaitGroup) {
			defer w.Done()
			solveResponse := board.Solve()
			results.AddSolution(index, solveResponse)

		}(idx, results, b, &wg)
	}

	wg.Wait()
	return len(boards)
}

func main() {

	boards, parseErr := solver.ParseFile("./data/top95.txt")
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	begin := time.Now()
	total := 0
	workers := runtime.NumCPU()

	fmt.Printf("Number of workers: %d\n", workers)
	fmt.Printf("Number of boards parsed %d\n", len(boards))

	result := &Results{
		Solutions: make(map[int]*solver.SolveResponse),
		Mutex:     sync.Mutex{},
	}

	batch := make(map[int]*solver.Board)
	for idx, board := range boards {
		if len(batch) == workers {
			// Let's solve the current batch
			solved := SolveBoardsInBatch(result, batch)
			total += solved
			// Clearing the buffer
			clear(batch)
			fmt.Printf("%d of %d boards have been processed\n", total, len(boards))
			// Adding the current element
			batch[idx] = board
		} else {
			batch[idx] = board
		}
	}

	// Checking the remaining
	if len(batch) > 0 {
		solved := SolveBoardsInBatch(result, batch)
		total += solved
	}

	// Printing all results
	result.PrintAll()

	fmt.Printf("Number of unsolved boards %d\n", result.NumberOfUnsolved())
	fmt.Printf("%d boards have been solved in parallel in %.2f seconds\n", total, time.Since(begin).Seconds())
}
