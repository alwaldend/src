
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
	remainder, rotation_start := make([]int, steps), length-steps
	// copy everything that needs to be shifted to another array
	copy(remainder, numbers[rotation_start:])
	// move everything to the right
	copy(numbers[steps:], numbers[0:rotation_start])
	// move shifted elements to the beginning
	copy(numbers[0:steps], remainder)
}
