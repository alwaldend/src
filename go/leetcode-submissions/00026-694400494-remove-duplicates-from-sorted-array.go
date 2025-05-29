func removeDuplicates(numbers []int) int {
	// ensure there are at least two numbers
	length := len(numbers)
	if length == 1 {
		return 1
	}
	index_non_duplicate := 1
	for index := 1; index < length; index++ {
		current, previous := numbers[index], numbers[index-1]
		if current == previous {
			// it is a duplicate - ignore it
			continue
		}
		// it is not a duplicate -> place it and move the index
		numbers[index_non_duplicate] = current
		index_non_duplicate++
	}
	return index_non_duplicate
}
