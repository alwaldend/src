func maxDepth(s string) int {
    maxDepth := 0
    curDepth := 0
    for _, ch := range s {
        if ch == '(' {
            curDepth++
            maxDepth = max(maxDepth, curDepth)
        } else if ch == ')' {
            curDepth--
        }
    }
    return maxDepth
}