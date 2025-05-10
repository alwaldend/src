func customSortString(order string, s string) string {
    count := make([]int, 26)

    for _, c := range s {
        count[c-'a']++
    }

    var result strings.Builder

    for _, c := range order {
        result.WriteString(strings.Repeat(string(c), count[c-'a']))
        count[c-'a'] = 0
    }

    for i := 0; i < 26; i++ {
        result.WriteString(strings.Repeat(string('a'+i), count[i]))
    }
    // UPVOTE :)
    return result.String()
}

