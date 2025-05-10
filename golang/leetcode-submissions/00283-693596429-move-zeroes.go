
func moveZeroes(numbers []int) {
	// ensure there are at least two numbers
	length := len(numbers)
	if length == 1 {
		return
	}
	result, index_result := make([]int, length), 0
	for _, number := range numbers {
		if number == 0 {
			continue
		}
		result[index_result] = number
		index_result++
	}
	copy(numbers, result)
}