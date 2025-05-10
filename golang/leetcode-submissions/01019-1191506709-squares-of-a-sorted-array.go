func sortedSquares(nums []int) []int {
    for i, num := range nums {
        if num < 0 {
            nums[i] = -num
        } else { 
            break
        }
    }
    slices.Sort(nums)
    for i, num := range nums {
        nums[i] = num * num
    }
    return nums
}