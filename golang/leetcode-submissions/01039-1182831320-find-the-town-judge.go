func findJudge(n int, trust [][]int) int {
    fromTo := make([][]int, n)
    toFrom := make([][]int, n)
    for _, trustArray := range trust {
        from, to := trustArray[0], trustArray[1]
        fromTo[from-1] = append(fromTo[from-1], to)
        toFrom[to-1] = append(toFrom[to-1], from)
    }
    for i, from := range toFrom {
        if len(from) == n - 1 && len(fromTo[i]) == 0 {
            return i + 1
        }
    }
    return -1
}