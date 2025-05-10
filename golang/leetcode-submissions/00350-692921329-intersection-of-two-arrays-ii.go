
func intersect(numbers_1 []int, numbers_2 []int) (result []int) {
	// 0 <= nums1[i], nums2[i] <= 1000
	count := make([]int, 1001)
	// counting how many times a number occured in the first array
	for _, number := range numbers_1 {
		count[number]++
	}
	for _, number := range numbers_2 {
		// the number did not occur in the first array
		// -> ignoring it
		if count[number] <= 0 {
			continue
		}
		count[number]--
		result = append(result, number)
	}
	return
}