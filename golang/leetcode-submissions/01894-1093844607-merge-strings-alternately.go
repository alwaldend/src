func mergeAlternately(word1 string, word2 string) string {
    var sb strings.Builder
    length1, length2 := len(word1), len(word2)
    for i := 0; i < max(length1, length2); i++ {
        if i == length1 {
            sb.WriteString(word2[i:length2])
            break
        }
        if i == length2 {
            sb.WriteString(word1[i:length1])
            break
        }
        sb.WriteByte(word1[i])
        sb.WriteByte(word2[i])
    }
    return sb.String()
}