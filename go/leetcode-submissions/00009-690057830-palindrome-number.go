
func isPalindrome(number int) (result bool) {
	if number < 0 {
		return false
	} else if number/10 == 0 {
		return true
	}
	current := number
	digits := []int{}
	for {
		if current == 0 {
			break
		}
		digit := current % 10
		current /= 10
		digits = append(digits, digit)
	}
	digits_length := len(digits)
	skip := -1
	if digits_length%2 != 0 {
		skip = digits_length / 2
	}
	for index, digit := range digits {
		if index != skip && digit != digits[digits_length-index-1] {
			return false
		}
	}
	return true
}
