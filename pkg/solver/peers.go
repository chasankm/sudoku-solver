package solver

var peerIDs [BoardSize * BoardSize][]int

func init() {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			id := cellID(row, col)
			seen := make(map[int]struct{}, 20)
			for i := 0; i < BoardSize; i++ {
				if i != col {
					seen[cellID(row, i)] = struct{}{}
				}
				if i != row {
					seen[cellID(i, col)] = struct{}{}
				}
			}
			rowMin := (row / BlockSize) * BlockSize
			colMin := (col / BlockSize) * BlockSize
			for r := rowMin; r < rowMin+BlockSize; r++ {
				for c := colMin; c < colMin+BlockSize; c++ {
					peer := cellID(r, c)
					if peer != id {
						seen[peer] = struct{}{}
					}
				}
			}
			peers := make([]int, 0, len(seen))
			for peer := range seen {
				peers = append(peers, peer)
			}
			peerIDs[id] = peers
		}
	}
}

func cellID(row int, col int) int {
	return row*BoardSize + col
}

func cellByID(data [BoardSize][BoardSize]*Cell, id int) *Cell {
	return data[id/BoardSize][id%BoardSize]
}

func candidateSetForPosition(data [BoardSize][BoardSize]*Cell, row int, col int) CandidateSet {
	digits := Digits.Clone()
	for _, peerID := range peerIDs[cellID(row, col)] {
		peer := cellByID(data, peerID)
		if peer.IsSolved() {
			digits = digits.AndNot(CandidateSetOf(int(peer.Value)))
		}
	}
	return digits
}
