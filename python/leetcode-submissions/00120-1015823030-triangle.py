class Solution:
    def minimumTotal(self, triangle: List[List[int]]) -> int:
        row_count = len(triangle)


        @cache
        def dp(row: int, column: int) -> int:
            if row == row_count or column == row + 1:
                return 0
        
            return triangle[row][column] + min(dp(row + 1, column), dp(row + 1, column + 1))

        return dp(0, 0)