func makeGood(s string) string {
	stack := make([]byte, 0, len(s))
	stack = append(stack, s[0])
	for i := 1; i < len(s); i++ {
		if len(stack) > 0 && getDiff(stack[len(stack)-1], s[i]) == 32 {
			stack = stack[:len(stack)-1]
			continue
		}
		stack = append(stack, s[i])
	}
	return string(stack)
}

func getDiff(a, b uint8) uint8 {
	if b > a {
		return b - a
	}
	return a - b
}