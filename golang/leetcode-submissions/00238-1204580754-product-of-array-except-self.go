func productExceptSelf(nums []int) []int {
    n := len(nums)
    res := make([]int, n)

    preProduct := 1
    for i := 0; i < n; i++ {
        res[i] = preProduct
        preProduct *= nums[i]
    }

    sufProduct := 1
    for i := n - 1; i >= 0; i-- {
        res[i] *= sufProduct
        sufProduct *= nums[i]
    }

    return res
}