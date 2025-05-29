func subarraysWithAtMostKDistinct(nums []int, k int) int {
    if k == 0 {
        return 0
    }

    countOccurrence := make(map[int]int)
    differentIntegers := 0
    left := 0
    result := 0

    for right := 0; right < len(nums); right++ {
        countOccurrence[nums[right]]++
        if countOccurrence[nums[right]] == 1 {
            differentIntegers++
        }

        for differentIntegers > k {
            countOccurrence[nums[left]]--
            if countOccurrence[nums[left]] == 0 {
                differentIntegers--
            }
            left++
        }

        result += right - left + 1
    }
    return result
}

func subarraysWithKDistinct(nums []int, k int) int {
    return subarraysWithAtMostKDistinct(nums, k) - subarraysWithAtMostKDistinct(nums, k-1)
}