func eliminateMaximum(dist []int, speed []int) int {
    arrival := []float32{}
    length := len(dist)
    for i := 0; i < length; i++ {
        arrival = append(arrival, float32(dist[i]) / float32(speed[i]))
    }
    slices.Sort(arrival)
    ans := 0
    for i := 0; i < length; i++ {
        if arrival[i] <= float32(i) {
            break
        }
        ans += 1
    }
    return ans
}