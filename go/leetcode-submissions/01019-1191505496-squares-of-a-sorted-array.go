func sortedSquares(nums []int) []int {
    for i, num := range nums {
        nums[i] = num * num
    }
    slices.Sort(nums)
    return nums
}