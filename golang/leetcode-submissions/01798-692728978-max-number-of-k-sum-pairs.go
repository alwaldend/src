
func maxOperations(numbers []int, sum int) int {
	count, previous := 0, make(map[int]int, len(numbers))
	for _, current := range numbers {
		target := sum - current
		target_unmatched, target_occured := previous[target]
		// if the current number has not occured before, then it is not in the map
		// -> it needs to be initialized
		if _, current_occured := previous[current]; !current_occured {
			previous[current] = 0
		}
		// number of duplicates of the current number has increased
		previous[current]++
		switch {
		case target_occured && target_unmatched == 0:
			// the target has occured before but there are no unmatched duplicates
			fallthrough
		case !target_occured:
			// in order to get the sum we need the target, but it has not appeared yet
			continue
		case target_occured && target_unmatched > 0:
			// the target has appeared before and there are some unmached duplicates left
			// -> removing the current number and the target from possible matches
			previous[current]--
			previous[target]--
			count++
		}
	}
	return count
}
