func findFarmland(land [][]int) [][]int {
	result := [][]int{}
	m, n := len(land), len(land[0])
	findFarmlandCoordinates := func(row, col int) []int {
		coordinates := []int{row, col}
		r, c := row, col
		for r < m && land[r][col] == 1 {
			r++
		}
		for c < n && land[row][c] == 1 {
			c++
		}
		coordinates = append(coordinates, r-1, c-1)
		for i := row; i < r; i++ {
			for j := col; j < c; j++ {
				land[i][j] = 0
			}
		}
		return coordinates
    }
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if land[i][j] == 1 {
				result = append(result, findFarmlandCoordinates(i, j))
			}
		}
	}
	return result
}
