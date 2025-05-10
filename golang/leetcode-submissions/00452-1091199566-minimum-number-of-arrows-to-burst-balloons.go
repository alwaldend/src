func findMinArrowShots(points [][]int) int {
	// greedy solution
	sort.Slice(points, func(i, j int) bool {
		return points[i][1] < points[j][1]
	})
	count := 1
	end := points[0][1]
	for i := 1; i < len(points); i++ {
		if points[i][0] > end {
			count++
			end = points[i][1]
		}
	}
	return count
}