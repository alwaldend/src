
func largestPerimeter(numbers []int) int {
	sort.Ints(numbers)
	for index := len(numbers) - 1; index > 1; index-- {
		current, sum_previous := numbers[index], numbers[index-1]+numbers[index-2]
		if current >= sum_previous {
			continue
		}
		return current + sum_previous
	}
	return 0
}