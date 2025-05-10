func countHomogenous(s string) int {
    var (
        mod int64 = 1000000007 
        total int64 = 0
        count int64 = 0
        cur = s[0]
    )
    
    for i := 0; i < len(s); i++ {
        char := s[i]
        if char == cur {
            count++
        } else {
            count = 1
            cur = char
        }
        total += count
    }
    
    return int(total % mod)
}