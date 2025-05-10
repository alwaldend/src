func findMaxLength(nums []int) int {
    hashmap := make(map[int]int)
    zeros, ones, maxLen := 0, 0, 0
    hashmap[0] = -1 
    for i, num := range nums {
        if num == 0 {
            zeros++
        } else {
            ones++
        }
        diff := zeros - ones
        if val, ok := hashmap[diff]; ok {
            maxLen = max(maxLen, i - val)
        } else {
            hashmap[diff] = i
        }
    }
    return maxLen
}