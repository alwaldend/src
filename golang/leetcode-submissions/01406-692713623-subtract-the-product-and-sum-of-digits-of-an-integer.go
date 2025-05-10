func subtractProductAndSum(number int) int {
	fist_digit := number % 10
	factorial, sum := fist_digit, fist_digit
	number /= 10
	for number > 0 {
		digit := number % 10
		factorial *= digit
		sum += digit
		number /= 10
	}
	return factorial - sum
}
