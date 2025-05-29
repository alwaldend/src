func maximumOddBinaryNumber(s string) string {
	l := len(s)
	res := make([]byte, l)
	index0, index1 := l-2, l-1
	for i := 0; i < l; i++ {
		if s[i] == '0' {
			res[index0] = '0'
			index0--
		} else {
			res[index1] = '1'
			index1 = (index1 + 1) % l
		}
	}
	return string(res)
}
