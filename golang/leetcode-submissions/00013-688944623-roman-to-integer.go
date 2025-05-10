
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
	types := dict{
		'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000,
	}
	for index, character := range roman {
		result += types[character]
        if index < 1 {
            continue}
        switch {
		case roman[index-1] == 'I' && in(character, "VX"):
			fallthrough
		case roman[index-1] == 'X' && in(character, "LC"):
			fallthrough
		case roman[index-1] == 'C' && in(character, "DM"):
			result -= types[rune(roman[index-1])] * 2
		}
	}
	return result
}

