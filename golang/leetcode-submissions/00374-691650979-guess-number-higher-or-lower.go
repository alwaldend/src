
func guessNumber(n int) int {
	// checking edge cases
	if guess(1) == 0 {
		return 1
	} else if guess(n) == 0 {
		return n
	}
	left, right := 1, n
	for right >= left {
		// overflow protection
		number := left + (right-left)/2
		switch guess(number) {
		case 0:
			// found the target
			return number
		case -1:
			// the number is bigger -> the target is to the left -> discard right
			right = number - 1
		case 1:
			// the number is smaller -> the target is to the right -> discard left
			left = number + 1
		}
	}
	// ide shows an error, this return is unreachable in this issue
	return 0
}
