func countSubarrays(nums []int, k int) int64 {
	maxValue := 0
	var maxValueIds []int
	var ans int64

	for i, x := range nums {
		if x > maxValue {
			maxValue, ans, maxValueIds = x, 0, []int{}
		}

		if x == maxValue {
			maxValueIds = append(maxValueIds, i)
		}

		if len(maxValueIds) >= k {
			ans += int64(maxValueIds[len(maxValueIds)-k]) + 1
		}
	}

	return ans
}