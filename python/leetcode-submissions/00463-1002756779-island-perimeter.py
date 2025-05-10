class Solution:
    def islandPerimeter(self, grid: List[List[int]]) -> int:
        perimeter = 0
        last_row = len(grid) - 1        
        last_cell = len(grid[0]) - 1
        for i, row in enumerate(grid):
            for j, cell in enumerate(row):
                if cell == 0:
                    continue
                
                if j == 0 or row[j-1] == 0:
                    perimeter += 1
                
                if j == last_cell or row[j+1] == 0:
                    perimeter += 1
                
                if i == 0 or grid[i-1][j] == 0:
                    perimeter += 1

                if i == last_row or grid[i+1][j] == 0:
                    perimeter += 1
        
        return perimeter
                