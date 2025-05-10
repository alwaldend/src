func findDuplicate(nums []int) int {
    seen := make(map[int]struct{}, len(nums) - 1)
    for _, num := range nums {
        if _, ok := seen[num]; ok {
            return num
        }
        seen[num] = struct{}{}
    }
    return -1
}