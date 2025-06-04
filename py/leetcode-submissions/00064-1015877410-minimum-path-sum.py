class Solution:
    def minPathSum(self, grid: List[List[int]]) -> int:
        row_count = len(grid)
        column_count = len(grid[0]) if grid else None
        if not row_count or not column_count:
            return 0

        dp = [[0] * column_count for _ in range(row_count)]
        dp[0][0] = grid[0][0]

        for column in range(1, column_count):
            dp[0][column] = grid[0][column] + dp[0][column - 1]

        for row in range(1, row_count):
            dp[row][0] = grid[row][0] + dp[row - 1][0]

            for column in range(1, column_count):
                dp[row][column] = grid[row][column] + min(
                    dp[row - 1][column], dp[row][column-1]
                )

        # @cache
        # def dp(row: int, column: int) -> int:
        #     if row == row_count or column == column_count:
        #         return 0
        #     return grid[row][column] + min(dp(row + 1, column), dp(row, column + 1))
        #return dp(0, 0)

        return dp[-1][-1]