func longestIdealString(s string, k int) int {
	dp := make([]int, 27)
	n := len(s)

	for i := n - 1; i >= 0; i-- {
		cc := s[i]
		idx := int(cc - 'a')
		maxi := -1 << 31

		left := max(idx-k, 0)
		right := min(idx+k, 26)

		for j := left; j <= right; j++ {
			maxi = max(maxi, dp[j])
		}

		dp[idx] = maxi + 1
	}

	max := -1 << 31
	for _, val := range dp {
		if val > max {
			max = val
		}
	}

	return max
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}