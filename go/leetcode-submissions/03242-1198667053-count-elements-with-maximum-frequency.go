func maxFrequencyElements(nums []int) int {
    mapCount := make(map[int]int)

    for _, val := range nums {
        mapCount[val]++
    }

    count := 0
    max := -1
    for _, freq := range mapCount {
        if freq > max {
            max = freq
        }
    }

    for _, freq := range mapCount {
        if freq == max {
            count += max
        }
    }
    return count
}