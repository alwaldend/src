func reversePrefix(word string, ch byte) string {
    var natija string
    for i := 0; i < len(word); i++ {
        if word[i] == ch { 
            ch_index := i
            for i >= 0 {
                natija += string(word[i])
                i--
            } 
            return natija + word[ch_index + 1 : ]
        }
    }    
    return word
}

