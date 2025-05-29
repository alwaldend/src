
func hammingWeight(number uint32) (result int) {
	for number > 0 {
		if number&1 != 0 {
			result++
		}
		number >>= 1
	}
	return
}
