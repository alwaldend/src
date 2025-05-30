func islandPerimeter(grid [][]int) int {
    islands := 0
    neighbors := 0

    for i := 0; i < len(grid); i++ {
        for j := 0; j < len(grid[0]); j++ {
            if grid[i][j] == 1 {
                islands++
                if i-1 >= 0 && grid[i-1][j] == 1 {
                    neighbors++
                }
                if j-1 >= 0 && grid[i][j-1] == 1 {
                    neighbors++
                }
            }
        }
    }

    return islands*4 - neighbors*2
}
