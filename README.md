## sudoku-solver

sudoku-solver is a simple application written in Go to solve any given sudoku board by using
different strategies to eliminate candidates and as a last resort back tracking if strategies 
do not produce any solutions anymore

The application uses the roaring bitmap https://github.com/RoaringBitmap/roaring for the set 
operations the rest of the code are pure golang code.

The solver algorithm uses sets of different strategies to eliminate the candidates and at some point
if strategies do not eliminate the candidates anymore the algorithm starts backtracking to find 
a unique valid solution by brute-forcing.

The strategies which are used in the algorithm

- Naked Quads strategy https://www.sudokuwiki.org/Naked_Candidates
- Naked Triples strategy https://www.sudokuwiki.org/Naked_Candidates
- Naked Pairs strategy https://www.sudokuwiki.org/Naked_Candidates
- XY Wings strategy https://www.learn-sudoku.com/xy-wing.html
- XYZ Wings strategy https://www.sudokuwiki.org/XYZ_Wing
- XWings strategy https://www.sudokuwiki.org/X_Wing_Strategy
- Sword Fish strategy https://www.sudokuwiki.org/Sword_Fish_Strategy
- Hidden Single strategy https://www.sudokuwiki.org/Hidden_Candidates
- Hidden Quads strategy https://www.sudokuwiki.org/Hidden_Candidates
- Hidden Triplets strategy https://www.sudokuwiki.org/Hidden_Candidates
- Hidden Pairs strategy https://www.sudokuwiki.org/Hidden_Candidates

The application parses two different board format which can be found on the data folder.

Example usage; for instance if you parse easy50.txt file application solves
board one by one and sample output is as below

In very rare cases application can not find unique solution or some inconsistency
happens with the initial state of the board, in that cases the Solve method returns an error
with the explanation and current state of the board.


```
go run cmd/main.go
Index: 0
Difficulty: Evil
Givens: 17
Is Solved: true
BackTracking used: false
Strategies used: [Naked Pairs, Hidden Single, Hidden Pairs, Naked Quads, Naked Triples]
Duration: 6.14 seconds
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
*_______*_______*______*
| 8 2 5 | 4 3 7 | 1 6 9 
| 7 9 1 | 5 8 6 | 4 3 2 
| 3 4 6 | 9 1 2 | 7 5 8 
*_______*_______*______*
| 2 8 9 | 6 4 3 | 5 7 1 
| 5 7 3 | 2 9 1 | 6 8 4 
| 1 6 4 | 8 7 5 | 2 9 3 

Index: 1
Difficulty: Evil
Givens: 17
Is Solved: true
BackTracking used: false
Strategies used: [Naked Triples, Naked Pairs, Hidden Single, Hidden Pairs, Naked Quads]
Duration: 5.84 seconds
Initial state: 
*_______*_______*______*
| 5 2 _ | _ _ 6 | _ _ _ 
| _ _ _ | _ _ _ | 7 _ 1 
| 3 _ _ | _ _ _ | _ _ _ 
*_______*_______*______*
| _ _ _ | 4 _ _ | 8 _ _ 
| 6 _ _ | _ _ _ | _ 5 _ 
| _ _ _ | _ _ _ | _ _ _ 
*_______*_______*______*
| _ 4 1 | 8 _ _ | _ _ _ 
| _ _ _ | _ 3 _ | _ 2 _ 
| _ _ 8 | 7 _ _ | _ _ _ 

Solution: 
*_______*_______*______*
| 5 2 7 | 3 1 6 | 4 8 9 
| 8 9 6 | 5 4 2 | 7 3 1 
| 3 1 4 | 9 8 7 | 5 6 2 
*_______*_______*______*
| 1 7 2 | 4 5 3 | 8 9 6 
| 6 8 9 | 2 7 1 | 3 5 4 
| 4 5 3 | 6 9 8 | 2 1 7 
*_______*_______*______*
| 9 4 1 | 8 2 5 | 6 7 3 
| 7 6 5 | 1 3 4 | 9 2 8 
| 2 3 8 | 7 6 9 | 1 4 5 

Index: 2
Difficulty: Evil
Givens: 17
Is Solved: true
BackTracking used: false
Strategies used: [Naked Pairs, Hidden Single, Hidden Pairs, Naked Quads, Naked Triples]
Duration: 6.47 seconds
Initial state: 
*_______*_______*______*
| 6 _ _ | _ _ _ | 8 _ 3 
| _ 4 _ | 7 _ _ | _ _ _ 
| _ _ _ | _ _ _ | _ _ _ 
*_______*_______*______*
| _ _ _ | 5 _ 4 | _ 7 _ 
| 3 _ _ | 2 _ _ | _ _ _ 
| 1 _ 6 | _ _ _ | _ _ _ 
*_______*_______*______*
| _ 2 _ | _ _ _ | _ 5 _ 
| _ _ _ | _ 8 _ | 6 _ _ 
| _ _ _ | _ 1 _ | _ _ _ 

Solution: 
*_______*_______*______*
| 6 1 7 | 4 5 9 | 8 2 3 
| 2 4 8 | 7 3 6 | 9 1 5 
| 5 3 9 | 1 2 8 | 4 6 7 
*_______*_______*______*
| 9 8 2 | 5 6 4 | 3 7 1 
| 3 7 4 | 2 9 1 | 5 8 6 
| 1 5 6 | 8 7 3 | 2 9 4 
*_______*_______*______*
| 8 2 3 | 6 4 7 | 1 5 9 
| 7 9 1 | 3 8 5 | 6 4 2 
| 4 6 5 | 9 1 2 | 7 3 8 

```
