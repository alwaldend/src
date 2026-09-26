package generator

import "math/rand/v2"

// compactCells accepts only moves that increase the occupied neighbor-pair count.
func compactCells(cells []bool, columns int, factor float64, random *rand.Rand) {
	attempts := int(factor * float64(8*len(cells)))
	for range attempts {
		from, to := random.IntN(len(cells)), random.IntN(len(cells))
		if cells[from] == cells[to] {
			continue
		}
		if !cells[from] {
			from, to = to, from
		}
		before := occupiedNeighbors(cells, columns, from)
		cells[from] = false
		if occupiedNeighbors(cells, columns, to) > before {
			cells[to] = true
		} else {
			cells[from] = true
		}
	}
}
