func firstMissingPositive(nums []int) int {
    n := len(nums)

	for i := 0; i < n; i++ {
		if nums[i] <= 0 || nums[i] > n {
			nums[i] = n + 1
		}
	}

	for i := 0; i < n; i++ {
		val := abs(nums[i])
		if val >= 1 && val <= n {
			flagIndex := val - 1
			if nums[flagIndex] > 0 {
				nums[flagIndex] *= -1
			}
		}
	}

	for i := 1; i <= n; i++ {
		if nums[i - 1] > 0 {
			return i
		}
	}

	return n + 1
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}