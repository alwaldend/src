func nearestValidPoint(x int, y int, points [][]int) int {
	point_target, smallest_distance, smallest_index := []int{x, y}, math.MaxInt, -1
	for index, point := range points {
		// ignore the point if it is not valid
		if !valid(point_target, point) {
			continue
		}
		distance := distance(point_target, point)
		// ignore the point if its Manhattan distance is bigger
		if distance >= smallest_distance {
			continue
		}
		// the distance is smaller, update the index and the distance
		smallest_index, smallest_distance = index, distance
	}
	// there are no valid points
	if smallest_index == -1 {
		return -1
	}
	return smallest_index
}

func valid(point_target []int, point []int) bool {
	switch {
	case point_target[0] == point[0]:
		fallthrough
	case point_target[1] == point[1]:
		return true
	default:
		return false
	}
}

func distance(point_target []int, point []int) int {
	return abs(point_target[0]-point[0]) + abs(point_target[1]-point[1])
}

func abs(number int) int {
	if number < 0 {
		return number * -1
	}
	return number
}
