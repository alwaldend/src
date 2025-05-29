func checkValidString(s string) bool {
    open := 0
    openMax := 0
    for _, char := range s {
        switch char {
        case '(':
            open++
            openMax++
        case ')':
            open--
            openMax--
        default:
            open--
            openMax++
        }
        if openMax < 0 {
            return false
        }
        if open < 0 {
            open = 0
        }
    }
    return open == 0
}