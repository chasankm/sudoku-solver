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
[0]: Sudoku (32) (Backtrack: false) is solved in 0.02 seconds
Difficulty: Medium
Givens: 32
BackTracking used: false
*_______*_______*______*
| 4 8 3 | 9 2 1 | 6 5 7 
| 9 6 7 | 3 4 5 | 8 2 1 
| 2 5 1 | 8 7 6 | 4 9 3 
*_______*_______*______*
| 5 4 8 | 1 3 2 | 9 7 6 
| 7 2 9 | 5 6 4 | 1 3 8 
| 1 3 6 | 7 9 8 | 2 4 5 
*_______*_______*______*
| 3 7 2 | 6 8 9 | 5 1 4 
| 8 1 4 | 2 5 3 | 7 6 9 
| 6 9 5 | 4 1 7 | 3 8 2 

[1]: Sudoku (30) (Backtrack: false) is solved in 0.03 seconds
Difficulty: Medium
Givens: 30
BackTracking used: false
*_______*_______*______*
| 2 4 5 | 9 8 1 | 3 7 6 
| 1 6 9 | 2 7 3 | 5 8 4 
| 8 3 7 | 5 6 4 | 2 1 9 
*_______*_______*______*
| 9 7 6 | 1 2 5 | 4 3 8 
| 5 1 3 | 4 9 8 | 6 2 7 
| 4 8 2 | 7 3 6 | 9 5 1 
*_______*_______*______*
| 3 9 1 | 6 5 7 | 8 4 2 
| 7 2 8 | 3 4 9 | 1 6 5 
| 6 5 4 | 8 1 2 | 7 9 3 

[2]: Sudoku (28) (Backtrack: false) is solved in 0.05 seconds
Difficulty: Hard
Givens: 28
BackTracking used: false
*_______*_______*______*
| 4 6 2 | 8 3 1 | 9 5 7 
| 7 9 5 | 4 2 6 | 1 8 3 
| 3 8 1 | 7 9 5 | 4 2 6 
*_______*_______*______*
| 1 7 3 | 9 8 4 | 2 6 5 
| 6 5 9 | 3 1 2 | 7 4 8 
| 2 4 8 | 5 6 7 | 3 1 9 
*_______*_______*______*
| 9 2 6 | 1 7 8 | 5 3 4 
| 8 3 4 | 2 5 9 | 6 7 1 
| 5 1 7 | 6 4 3 | 8 9 2 
```
