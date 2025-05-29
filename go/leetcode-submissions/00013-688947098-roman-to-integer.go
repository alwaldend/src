func romanToInt(input string) int {
	var result int
	types := map[rune]int{
		'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000,
	}
	var previous rune
	for _, character := range input {
		result += types[character]
		switch {
		case previous == 'I' && (character == 'V' || character == 'X'):
			fallthrough
		case previous == 'X' && (character == 'L' || character == 'C'):
			fallthrough
		case previous == 'C' && (character == 'D' || character == 'M'):
			result -= types[previous] * 2
		}
		previous = character
	}
	return result
}
