func insert(intervals [][]int, new []int) [][]int {
    n := len(intervals)
    i := sort.Search(n, func(i int) bool { return intervals[i][0] > new[0] })
    j := sort.Search(n, func(j int) bool { return intervals[j][1] > new[1] })
    if i >= 1 && new[0] <= intervals[i-1][1] {
        new[0] = intervals[i-1][0]
        i--
    }
    if j < n && new[1] >= intervals[j][0] {
        new[1] = intervals[j][1]
        j++
    }
    return append(intervals[:i], append([][]int{new}, intervals[j:]...)...)
}