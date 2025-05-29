func getWinner(arr []int, k int) int {
    curMax, winCount := arr[0], 0
    i, length := 1, len(arr)
    for {
        num := arr[i]
        if num > curMax {
            curMax = num
            winCount = 1
        } else {
            winCount++
        }
        if winCount == k {
            return curMax
        }
        i = (i + 1) % length
    }
    return -1
}