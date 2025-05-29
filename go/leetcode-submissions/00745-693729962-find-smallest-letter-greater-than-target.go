

func nextGreatestLetter(letters []byte, target byte) byte {
	length := len(letters)
	left, right := 0, length-1
	for right >= left {
		index := left + (right-left)/2
		character := letters[index]
		// the character is bigger, we found at least one result
		// smaller characters are to the left -> discard right
		if character > target {
			right = index - 1
			continue
		}
		// character is either equal or smaller -> there is no results to the
		// left -> discard left
		left = index + 1
	}
	return letters[left%length]
}
