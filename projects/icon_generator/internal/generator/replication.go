package generator

import "math/rand/v2"

// Dimensions are limited to 4096 by 4096, so an index needs at most 24 bits.
// Store the three-bit heading above it to keep the queue bounded at 64 MiB.
const replicationIndexBits = 24

type replicationVacancy struct {
	index, direction, neighbors int
}

// replicateCells gives each cell one chance to add a pair of adjacent children.
// Parents remain visible. Every child occupies a vacancy and is queued once,
// bounding both the number of decisions and the queue by the canvas cell count.
func replicateCells(cells []bool, columns int, options Options, field *regionalField, random, headings *rand.Rand) {
	directed := options.ReplicationForwardBias > 0
	weighted := directed || options.ReplicationCrowding > 0 || options.ReplicationBranching > 0
	queue := make([]uint32, 0, len(cells))
	for index, occupied := range cells {
		if occupied {
			entry := uint32(index)
			if directed {
				entry |= uint32(headings.IntN(8)) << replicationIndexBits
			}
			queue = append(queue, entry)
		}
	}
	rows := len(cells) / columns
	for head := 0; head < len(queue); head++ {
		index := int(queue[head] & ((1 << replicationIndexBits) - 1))
		x, y := index%columns, index/columns
		probability := options.Clusterization
		if field != nil && options.ReplicationGrowthVariation > 0 {
			probability = field.growthProbability(x, y)
		}
		if random.Float64() >= probability {
			continue
		}
		heading := int(queue[head] >> replicationIndexBits)
		var vacancies [8]replicationVacancy
		count := 0
		for direction, offset := range neighborDirections {
			nx, ny := x+offset.X, y+offset.Y
			if nx < 0 || nx >= columns || ny < 0 || ny >= rows || cells[ny*columns+nx] {
				continue
			}
			vacancy := replicationVacancy{index: ny*columns + nx, direction: direction}
			if options.ReplicationCrowding > 0 {
				// Every vacancy is adjacent to the parent; exclude that contact.
				vacancy.neighbors = occupiedNeighbors(cells, columns, vacancy.index) - 1
			}
			vacancies[count] = vacancy
			count++
		}
		if count < 2 {
			continue
		}
		var first, second int
		if weighted {
			first = replicationChoice(vacancies[:count], heading, options.ReplicationForwardBias, options.ReplicationCrowding, -1, random)
			secondHeading, secondBias := heading, options.ReplicationForwardBias
			if options.ReplicationBranching > 0 && random.Float64() < options.ReplicationBranching {
				secondHeading = (vacancies[first].direction + 2 + 4*random.IntN(2)) % 8
				secondBias = 1
			}
			second = replicationChoice(vacancies[:count], secondHeading, secondBias, options.ReplicationCrowding, first, random)
		} else {
			// Preserve the original random stream exactly when controls are off.
			first, second = random.IntN(count), random.IntN(count-1)
			if second >= first {
				second++
			}
		}
		for _, child := range [...]replicationVacancy{vacancies[first], vacancies[second]} {
			cells[child.index] = true
			entry := uint32(child.index)
			if directed {
				direction := child.direction
				if headings.Float64() >= options.ReplicationInheritance {
					direction = headings.IntN(8)
				}
				entry |= uint32(direction) << replicationIndexBits
			}
			queue = append(queue, entry)
		}
	}
}

// replicationChoice keeps all eligible weights positive. Preferences never
// prevent a birth when two vacancies exist. Exclude is the first child's slot;
// account for its new contact while scoring the second child.
func replicationChoice(vacancies []replicationVacancy, heading int, bias, crowding float64, exclude int, random *rand.Rand) int {
	var weights [8]float64
	total, last := 0.0, 0
	for slot, vacancy := range vacancies {
		if slot == exclude {
			continue
		}
		turn := (vacancy.direction - heading + 8) % 8
		turn = min(turn, 8-turn)
		preference := [...]float64{16, 8, 2, 0.5, 0.125}[turn]
		weight := 1 + bias*(preference-1)
		neighbors := vacancy.neighbors
		if crowding > 0 && exclude >= 0 {
			first, second := neighborDirections[vacancies[exclude].direction], neighborDirections[vacancy.direction]
			dx, dy := first.X-second.X, first.Y-second.Y
			if dx >= -1 && dx <= 1 && dy >= -1 && dy <= 1 {
				neighbors++
			}
		}
		weight /= 1 + 4*crowding*float64(neighbors*neighbors)
		weights[slot] = weight
		total += weight
		last = slot
	}
	choice := random.Float64() * total
	for slot := range vacancies {
		choice -= weights[slot]
		if weights[slot] > 0 && choice < 0 {
			return slot
		}
	}
	return last // Floating-point rounding at the top of the cumulative range.
}
