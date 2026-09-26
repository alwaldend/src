package generator

import (
	"math"
	"math/rand/v2"
)

type growingTip struct {
	x, y, heading, curvature float64
	recent                   [3]int
	junction, remaining, age int
}

// growTips reserves the initial population before growing shuffled roots.
// Only root indices span the canvas; floating-point state belongs to one root's
// frontier. Forks divide a finite budget rather than multiplying it.
func growTips(cells []bool, columns, initial int, factor float64, random *rand.Rand) {
	roots := make([]uint32, 0, initial)
	for index, occupied := range cells {
		if occupied {
			roots = append(roots, uint32(index))
		}
	}
	random.Shuffle(len(roots), func(i, j int) { roots[i], roots[j] = roots[j], roots[i] })
	maxLength := 8 + int(1024*factor*factor)
	var pending []growingTip
	for _, root := range roots {
		index := int(root)
		current := growingTip{
			x: float64(index%columns) + 0.5, y: float64(index/columns) + 0.5,
			heading: random.Float64() * 2 * math.Pi, curvature: (random.Float64() - 0.5) * 0.05,
			recent: [3]int{index, -1, -1}, junction: -1,
			remaining: maxLength/2 + random.IntN(maxLength-maxLength/2+1),
		}
		pending = append(pending[:0], current)
		for len(pending) > 0 {
			current = pending[len(pending)-1]
			pending = pending[:len(pending)-1]
			for current.remaining > 0 && advanceTip(cells, columns, &current, random) {
				if current.age < 12 || current.remaining < 24 || random.Float64() >= 0.025 {
					continue
				}
				branch := current
				branch.junction = current.recent[0]
				branch.age = 0
				branch.remaining = int(float64(current.remaining) * (0.2 + 0.25*random.Float64()))
				allocation := branch.remaining
				branch.heading += float64(2*random.IntN(2)-1) * (1.0 + 0.4*random.Float64())
				// Reserve the first side cell now; deferred growth cannot turn
				// the branch into a disconnected path or overwrite the main tip.
				if advanceTip(cells, columns, &branch, random) {
					current.remaining -= allocation
					current.age = 0
					current.junction = branch.junction
					pending = append(pending, branch)
				}
			}
		}
	}
}

// advanceTip rasterizes a continuous half-cell walk into one adjacent cell.
// Three substeps suffice to leave a unit square along a fixed heading, and the
// explicit limit also bounds rounding edge cases. Failed steps end this tip.
func advanceTip(cells []bool, columns int, tip *growingTip, random *rand.Rand) bool {
	tip.curvature = 0.96*tip.curvature + 0.006*(2*random.Float64()-1)
	tip.heading += tip.curvature
	dy, dx := math.Sincos(tip.heading)
	rows := len(cells) / columns
	for range 3 {
		tip.x += dx * 0.5
		tip.y += dy * 0.5
		x, y := int(math.Floor(tip.x)), int(math.Floor(tip.y))
		if x < 0 || x >= columns || y < 0 || y >= rows {
			return false
		}
		index := y*columns + x
		if index == tip.recent[0] {
			continue
		}
		if cells[index] || tipContact(cells, columns, index, tip) {
			return false
		}
		cells[index] = true
		tip.recent = [3]int{index, tip.recent[0], tip.recent[1]}
		tip.remaining--
		tip.age++
		return true
	}
	return false
}

// Keep a vacant cell between unrelated paths. Recent trail cells and the small
// fork junction are exempt so a rasterized curve can turn and forks can depart.
func tipContact(cells []bool, columns, index int, tip *growingTip) bool {
	x, y := index%columns, index/columns
	for _, offset := range neighborDirections {
		nx, ny := x+offset.X, y+offset.Y
		if nx < 0 || nx >= columns || ny < 0 || ny >= len(cells)/columns {
			continue
		}
		neighbor := ny*columns + nx
		if !cells[neighbor] || neighbor == tip.recent[0] || neighbor == tip.recent[1] || neighbor == tip.recent[2] {
			continue
		}
		if tip.junction >= 0 {
			jx, jy := tip.junction%columns, tip.junction/columns
			if nx >= jx-2 && nx <= jx+2 && ny >= jy-2 && ny <= jy+2 {
				continue
			}
		}
		return true
	}
	return false
}
