func minOperations(nums []int, k int) int {
    res := k
    for _, n := range nums {
        res ^= n
    }
    var ans int
    for res > 0 {
        ans += res%2
        res = res >> 1
    }
    return ans
}

