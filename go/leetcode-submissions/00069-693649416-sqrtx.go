
func mySqrt(number int) int {
	left, right := 0, number
	for right >= left {
		current := left + (right-left)/2
		square_current := current * current
		square_next := (current + 1) * (current + 1)
		switch {
		case square_current <= number && square_next > number:
			// found the target
			return current
		case square_current > number:
			// target is to the left -> discard right
			right = current - 1
		case square_current < number:
			// target is to the right -> discard left
			left = current + 1
		}
	}
	return -1
}