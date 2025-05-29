func findArray(pref []int) []int {
    prev := pref[0]
    for i := 1; i < len(pref); i++ {
        prev, pref[i] = pref[i], prev ^ pref[i] 
    }
    return pref
}