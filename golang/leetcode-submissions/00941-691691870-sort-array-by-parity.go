
func sortArrayByParity(numbers []int) []int {
	length := len(numbers)
	if length == 1 {
		return numbers
	}
	results, index_even, index_odd := make([]int, length), 0, length-1
	for _, number := range numbers {
		if (number & 1) == 0 {
			// number is even -> place it from the beginning
			results[index_even] = number
			index_even++
		} else {
			// number is odd -> place it from the end
			results[index_odd] = number
			index_odd--
		}
	}
	return results
}
