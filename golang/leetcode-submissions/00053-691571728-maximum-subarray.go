
func maxSubArray(numbers []int) int {
	length := len(numbers)
	switch length {
	case 0:
		return 0
	case 1:
		return numbers[0]
	}
	maxCurrent, maxOverall := 0, math.MinInt
	for _, number := range numbers {
		maxCurrent += number
		if maxCurrent > maxOverall {
			maxOverall = maxCurrent
		}
		if maxCurrent < 0 {
			maxCurrent = 0
		}
	}
	return maxOverall
}