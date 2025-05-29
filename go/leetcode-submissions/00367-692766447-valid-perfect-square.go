
func isPerfectSquare(number int) bool {
	left, right := 1, number
	for right >= left {
		current := left + (right-left)/2
		square := current * current
		switch {
		case square == number:
			// found the target
			return true
		case square > number:
			// square is bigger -> root is to the left -> discard right
			right = current - 1
		case square < number:
			// square is smaller -> root is to the right -> discard left
			left = current + 1
		}
	}
	return false
}