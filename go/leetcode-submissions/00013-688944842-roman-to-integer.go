
func in(character rune, targets string) bool {
	for _, target := range targets {
		if character == target {
			return true
		}
	}
	return false
}

func romanToInt(roman string) int {
	type dict map[rune]int
	var result int
	input := []rune(roman)
	types := dict{
		'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000,
	}
	var previous rune
	for _, character := range input {
		result += types[character]
		switch {
		case previous == 'I' && in(character, "VX"):
			fallthrough
		case previous == 'X' && in(character, "LC"):
			fallthrough
		case previous == 'C' && in(character, "DM"):
			result -= types[previous] * 2
		}
		previous = character
	}
	return result
}
