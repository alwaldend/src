func maxSubarrayLength(nums []int, k int) int {
    i := 0
    j := 0
    n := len(nums)
    ans := 1
    mp := make(map[int]int)

    for i < n {
        mp[nums[i]]++
        for mp[nums[i]] > k {
            mp[nums[j]]--
            j++
        }
        if i-j+1 > ans {
            ans = i - j + 1
        }
        i++
    }
    return ans
}