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

func maxSubarraySumCircular(nums []int) int {
    var (
        max_sum_ending_here = 0
        min_sum_ending_here = 0
        max_sum_so_far = math.MinInt
        min_sum_so_far = math.MaxInt
        total = 0
    )
    for i := 0; i < len(nums); i++ {
        num := nums[i]
        total += num
        max_sum_ending_here += num
        min_sum_ending_here += num
        max_sum_so_far = max(max_sum_so_far, max_sum_ending_here)
        min_sum_so_far = min(min_sum_so_far, min_sum_ending_here)
        max_sum_ending_here = max(max_sum_ending_here, 0)
        min_sum_ending_here = min(min_sum_ending_here, 0)
    }
    if max_sum_so_far < 0 {
        return max_sum_so_far
    }
    return max(max_sum_so_far, total - min_sum_so_far)
}