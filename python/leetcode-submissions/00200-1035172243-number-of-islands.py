class Solution:
    def numIslands(self, grid: List[List[str]]) -> int:
        row_count, col_count = len(grid), len(grid[0])
        island_count = 0
        moves = ((0, 1), (0, -1), (1, 0), (-1, 0))
        queue = set((row, col) 
                    for row in range(row_count) 
                    for col in range(col_count) 
                    if grid[row][col] == "1")

        @cache
        def remove_island(row: int, col: int) -> None:
            if not 0 <= row < row_count or not 0 <= col < col_count:
                return

            for row_delta, col_delta in moves:
                new_cell = (row + row_delta, col + col_delta)
                if new_cell in queue:
                    queue.remove(new_cell)
                    remove_island(*new_cell)

        while queue:
            remove_island(*queue.pop())
            island_count += 1

        return island_count

            