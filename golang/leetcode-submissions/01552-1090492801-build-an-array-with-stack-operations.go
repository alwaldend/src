func buildArray(target []int, n int) []string {
    ops := []string{}
    length := len(target)
    matchNext, matchNextVal := 0, target[0]
    for i := 1; i <= n; i++ {
        if i == matchNextVal {
            ops = append(ops, "Push")
            matchNext += 1
            if matchNext == length {
                break
            }
            matchNextVal = target[matchNext]
        } else {
            ops = append(ops, "Push", "Pop")
        }
    }
    return ops
}