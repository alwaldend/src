func minWindow(s string, t string) string {
    m := len(s)
    n := len(t)
    
    if m == 0 || n == 0 || m < n {
        return ""
    }

    dict := make(map[byte]int)
    for i := 0; i < n; i++ {
        dict[t[i]]++
    }
    
    // required unique chars
    required := len(dict)
    actual := 0
    window := make(map[byte]int)
    minSize := math.MaxInt64
    start := 0
    left, right := 0, 0
    
    for end := 0; end < m; end++ {
        c := s[end]
        window[c]++
        
        if value, ok := dict[c]; ok {
            if value == window[c] {
                actual++
            }
        }
        
        for start <= end && actual == required {
            size := end-start+1
            if size < minSize {
                minSize = size
                left = start
                right = end
            }
            
            rc := s[start]
            window[rc]--
            if value, ok := dict[rc]; ok {
                if value > window[rc] {
                    actual--
                }
            }
            start++
        }
    }
    
    if minSize == math.MaxInt64 {
        return ""
    }
    return s[left:right+1]
}