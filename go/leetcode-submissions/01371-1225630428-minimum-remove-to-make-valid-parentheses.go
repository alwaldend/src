type stack struct {
    top   int
    v     []byte
    index []int
}

func minRemoveToMakeValid(s string) string {
    stc := stack{top: -1}

    // collect bad brackets
    for i, v := range s {
        if v == '(' {
            stc.v = append(stc.v, '(')
            stc.index = append(stc.index, i)
            stc.top++
        } else if v == ')' {
            if stc.top > -1 && stc.v[stc.top] == '(' {
                stc.v = stc.v[:stc.top]
                stc.index = stc.index[:stc.top]
                stc.top--
            } else {
                stc.v = append(stc.v, ')')
                stc.index = append(stc.index, i)
                stc.top++
            }
        }
    }

    // remove them
    res := []byte{}
    i := 0
    for _, v := range stc.index {
        res = append(res, s[i:v]...)
        i = v + 1
    }

    // checking of end
    if len(res) + len(stc.index) < len(s) {
        res = append(res, s[i:]...)
    }

    return string(res)
}