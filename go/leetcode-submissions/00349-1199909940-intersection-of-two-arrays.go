func intersection(nums1 []int, nums2 []int) []int {
    seen := make([]int, 1000)
    for i := range nums1 {
        seen[nums1[i]]++
    }

    res := make([]int, 0)
    for i := range nums2 {
        if seen[nums2[i]] > 0 {
            res = append(res, nums2[i])
            seen[nums2[i]] = 0
        }
    }

    return res
}