
func findTheDistanceValue(array_1 []int, array_2 []int, target int) int {
	// sorting it to use binary search
	sort.Ints(array_2)
	length_2, count := len(array_2), 0
	for _, current_1 := range array_1 {
		left, right, add_to_count := 0, length_2-1, true
		for right >= left {
			index := left + (right-left)/2
			current_2 := array_2[index]
			if abs(current_1, current_2) <= target {
				// current_1 is inside |arr1[i]-arr2[j]| <= d
				// -> ignore it
				add_to_count = false
				break
			}
			switch {
			case current_2 > current_1:
				// current_2 is bigger than current_1
				// -> all numbers to the right are bigger
				// -> discard right, add to count
				right = index - 1
			case current_2 < current_1:
				// current_2 is smaller than current_1
				// -> all numbers to the left are smaller
				// -> discard left, add to count
				left = index + 1
			}
		}
		if add_to_count {
			count++
		}
	}
	return count
}

func abs(number_1 int, number_2 int) int {
	difference := number_1 - number_2
	if difference < 0 {
		return difference * -1
	}
	return difference
}
