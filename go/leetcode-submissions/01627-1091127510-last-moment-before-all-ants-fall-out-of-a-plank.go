
func getLastMoment(n int, left []int, right []int) int {
    return max(slices.Max(append(left, 0)), n - slices.Min(append(right, n)))
}