func findMinArrowShots(points [][]int) int {
    // Sort the balloons based on their end coordinates
    sort.Slice(points, func(i, j int) bool {
        return points[i][1] < points[j][1]
    })

    arrows := 1
    prevEnd := points[0][1]

    // Count the number of non-overlapping intervals
    for i := 1; i < len(points); i++ {
        if points[i][0] > prevEnd {
            arrows++
            prevEnd = points[i][1]
        }
    }

    return arrows
}