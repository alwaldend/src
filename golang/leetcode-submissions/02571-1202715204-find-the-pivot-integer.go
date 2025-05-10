func getSum(x int) int {
	return x * (x + 1) / 2
}

func pivotInteger(n int) int {
    sum := getSum(n)

	l, r := 1, n
	for l <= r {
		m := (l + r) / 2
		firstPart := getSum(m)
		secondPart := sum - firstPart + m

		if firstPart == secondPart {
			return m
		} else if firstPart > secondPart {
			r = m - 1
		} else {
			l = m + 1
		}
	}

	return -1
}