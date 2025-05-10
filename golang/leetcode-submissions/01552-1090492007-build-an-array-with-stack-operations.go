func buildArray(target []int, n int) []string {
    ops := []string{}
    length := len(target)
    matchNext := 0
    for i := 1; i <= n; i++ {
        if matchNext == length {
            break
        }
        ops = append(ops, "Push")
        if i == target[matchNext] {
            matchNext += 1
        } else {
            ops = append(ops, "Pop")
        }
    }
    return ops
}