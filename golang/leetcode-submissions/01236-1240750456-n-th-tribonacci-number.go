func tribonacci(n int) int {
	if n < 2 {
		return n
	}

	dp := []int{0, 1, 1}

	for i := 3; i <= n; i++ {
		next := dp[0] + dp[1] + dp[2]
		dp[0], dp[1], dp[2] = dp[1], dp[2], next
	}

	return dp[2]
}