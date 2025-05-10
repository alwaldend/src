
func isValid(inputString string) bool {
	if len(inputString) == 0 || (len(inputString)&1) != 0 {
		return false
	}
	occurences := []rune{}
	ends := map[rune]rune{')': '(', '}': '{', ']': '['}
	for _, character := range inputString {
		length := len(occurences)
		last_valid, is_end := ends[character]
		if (is_end && length == 0) || (is_end && occurences[length-1] != last_valid) {
			return false
		}
		if is_end {
			occurences = occurences[0 : length-1]
		} else {
			occurences = append(occurences, character)
		}
	}
	if len(occurences) > 0 {
		return false
	}
	return true
}
