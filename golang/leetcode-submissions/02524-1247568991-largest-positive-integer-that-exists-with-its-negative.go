func findMaxK(nums []int) int {
    alreadySeenNums := map[int]bool{}
    res := -1

    for _, num := range nums {
        numAbs := int(math.Abs(float64(num)))

        if numAbs < res {
            continue
        }

        if alreadySeenNums[-num] {
            res = numAbs
        } else {
            alreadySeenNums[num] = true
        }
    }

    return res
}