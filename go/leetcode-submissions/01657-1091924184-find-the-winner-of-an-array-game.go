func getWinner(arr []int, k int) int {
    if k == 1 {
        return max(arr[0], arr[1])
    }
    length := len(arr)
    if k >= length {
        return slices.Max(arr)
    }
    curWinner, winCount := arr[0], 0
    for _, num := range arr[1:] {
        if curWinner > num {
            winCount++
        } else {
            curWinner = num
            winCount = 1
        }
        if winCount == k {
            return curWinner
        }
    }
    return curWinner
}