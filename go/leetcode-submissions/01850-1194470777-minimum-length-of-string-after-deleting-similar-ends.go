func minimumLength(s string) int {
    left := 0
    right := len(s) - 1
    
    for left < right {
        lc := s[left]
        rc := s[right]
        
        if lc != rc {
            return right - left + 1
        }
        
        for left + 1 < right && lc == s[left + 1] {
            left++
        }
        
        for left < right - 1 && rc == s[right - 1] {
            right--
        }
        
        right--
        left++
    }
    
    return right - left + 1
}