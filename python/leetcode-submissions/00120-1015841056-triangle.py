class Solution:
    def minimumTotal(self, triangle: List[List[int]]) -> int:
        row_count = len(triangle)
        dp = [[0] * row_count for _ in range(row_count)]
        dp[0][0] = triangle[0][0]

        for row in range(1, row_count):
            dp[row][0] = triangle[row][0] + dp[row-1][0]
            dp[row][row] = triangle[row][row] + dp[row-1][row-1]

            for column in range(1, row):
                dp[row][column] = triangle[row][column] + min(
                    dp[row-1][column], dp[row-1][column-1]
                )

        return min(dp[-1])