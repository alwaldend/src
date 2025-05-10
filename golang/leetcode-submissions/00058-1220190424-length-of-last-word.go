func lengthOfLastWord(s string) int {
    length := len(s)
    count := 0
    for i := length - 1; i >= 0; i-- {
        isSpace := s[i] == ' '
        if isSpace && count != 0 {
            return count
        } else if isSpace {
            continue
        }
        count++
    }
    return count
}