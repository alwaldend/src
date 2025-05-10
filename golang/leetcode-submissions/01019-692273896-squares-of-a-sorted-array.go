
func sortedSquares(numbers []int) []int {
	length := len(numbers)
	switch length {
	case 0:
		return numbers
	case 1:
		return []int{numbers[0] * numbers[0]}
	}
	index_left, index_right := 0, length-1
	results := make([]int, length)
	for index := index_right; index_left <= index_right; index-- {
		left_square := numbers[index_left] * numbers[index_left]
		right_square := numbers[index_right] * numbers[index_right]
		if right_square > left_square {
			results[index] = right_square
			index_right--
			continue
		}
		results[index] = left_square
		index_left++
	}
	return results
}
