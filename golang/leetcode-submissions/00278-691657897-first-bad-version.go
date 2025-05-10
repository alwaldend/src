
func firstBadVersion(n int) int {
	// checking edge cases
	if isBadVersion(1) {
		return 1
	}
	left, right, bad_version := 1, n, math.MaxInt
	for right >= left {
		// overflow protection
		version := left + (right-left)/2
		switch isBadVersion(version) {
		case false:
			// it is good -> versions to the left are good -> discard left
			left = version + 1
		case true:
			// it is bad -> the first bad version is either this one or to the left
			// discard right
			bad_version = version
			right = version - 1
		}
	}
	// the search space is empty
	return bad_version
}
