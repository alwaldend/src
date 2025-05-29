
func rotate(numbers []int, steps int) {
	length := len(numbers)
	// removing unnecessary steps
	if steps >= length {
		steps %= length
	}
	// checking edge cases
	if length == 1 || steps == 0 {
		return
	}
	results := make([]int, length)
	for index, number := range numbers {
		results[(index+steps)%length] = number
	}
	copy(numbers, results)
}