class Solution:
    def numIslands(self, grid: List[List[str]]) -> int:
        row_count, col_count = len(grid), len(grid[0])
        island_count = 0
        moves = ((0, 1), (0, -1), (1, 0), (-1, 0))

        @cache
        def dfs(row: int, col: int) -> None:
            if not 0 <= row < row_count or not 0 <= col < col_count or grid[row][col] == "0":
                return
            
            grid[row][col] = "0"
            for row_delta, col_delta in moves:
                dfs(row + row_delta, col + col_delta)

        for row in range(row_count):
            for col in range(col_count):
                if grid[row][col] == "0":
                    continue

                dfs(row, col)
                island_count += 1

        return island_count

            