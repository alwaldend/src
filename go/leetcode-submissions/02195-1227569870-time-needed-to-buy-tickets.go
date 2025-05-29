func timeRequiredToBuy(tickets []int, k int) int {
    n := len(tickets)
    d := tickets[k]
    res := 0
    for i := 0; i <= k; i++ {
        res += min(d, tickets[i])
    }
    for i := k + 1; i < n; i++ {
        res += min(d - 1, tickets[i])
    }
    return res
}