func getCommon(nums1 []int, nums2 []int) int {
    var l, r int
    for l < len(nums1) && r < len(nums2) {
        switch {
        case nums1[l] == nums2[r]:
            return nums1[l]
        case nums1[l] < nums2[r]:
            l++
        default:
            r++
        }
    }
    return -1
}