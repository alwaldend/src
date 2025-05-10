class Solution:
    def uniquePaths(self, m: int, n: int) -> int:

        @cache
        def dp(row: int, col: int) -> int:
            if not 0 <= row < m or not 0 <= col < n:
                return 0
            
            if row == m - 1 and col == n - 1:
                return 1
            
            return dp(row + 1, col) + dp(row, col + 1)

        return dp(0, 0)