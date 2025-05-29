func ladderLength(beginWord string, endWord string, wordList []string) int {
    
    set := make(map[string]struct{}, len(wordList))
    
    present := false
    
    for _, v := range wordList {
        if endWord == v {
            present = true
        } 
        
        set[v] = struct{}{} 
    }
    
    if !present {
        return 0
    }
    
    set[beginWord] = struct{}{}
    q := []string{beginWord}
    
    depth := 1
    breadth := 0
    
    breadth = len(q)
    
    for ;breadth > 0; {
        s := q[0]
        
        if s == endWord {
            return depth
        }
        
        for i:='a'; i <= 'z'; i += 1 {
                        
            for j := 0; j<len(s); j++ {
 
                if rune(s[j]) != i {
                    
                    temp := s[:j] + string(i) + s[j+1:]
                    if _, ok := set[temp]; !ok {
                        continue
                    }
                    
                    q = append(q, temp)
                    delete(set, s)
                }
            }
        }
        
        q = q[1:]
        
        breadth -= 1
        if breadth == 0 {
            breadth = len(q)
            depth += 1
        }
    }
    
    return 0 
}