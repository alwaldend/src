func numSubarrayProductLessThanK(nums []int, k int) int {
	if k <= 1 {
		return 0
	}

	l := 0
	p := 1
	res := 0
	for r, num := range nums {
		p *= num
		for p >= k {
			p /= nums[l]
			l++
		}
        
		res += r - l + 1
	}

	return res
}