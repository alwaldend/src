func calculate(s string) int {
	result, _ := calculateFrom(s, 0)
	return result
}

func calculateFrom(s string, idFrom int) (result, idEnd int) {
	result, currNum, sign := 0, 0, 1

	for idEnd = idFrom; idEnd < len(s) && s[idEnd] != ')'; idEnd++ {
		switch {
		case s[idEnd] >= '0':
			currNum = currNum*10 + int(s[idEnd]-'0')
		case s[idEnd] == '(':
			currNum, idEnd = calculateFrom(s, idEnd+1)
		case s[idEnd] == '-' || s[idEnd] == '+':
			result, currNum = result+currNum*sign, 0
			sign = 44 - int(s[idEnd]) // '-'=45; '+'=43
		}
	}

	return result + currNum*sign, idEnd
}