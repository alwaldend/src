func findArray(pref []int) []int {
    prev := pref[0]
    for i := 1; i < len(pref); i++ {
        cur := pref[i]
        prev, pref[i] = cur, prev ^ cur 
    }
    return pref
}