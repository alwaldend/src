class Solution:
    def uniquePathsWithObstacles(self, obstacleGrid: List[List[int]]) -> int:
        row_count = len(obstacleGrid)
        column_count = len(obstacleGrid[0])
        # dp = [[0] * column_count for _ in range(row_count)]

        # for row in range(row_count):
        #     for column in range(column_count):

        @cache
        def dp(row: int, column: int) -> int:
            if row == row_count or column == column_count or obstacleGrid[row][column] == 1:
                return 0
            
            if row == row_count - 1 and column == column_count - 1:
                return 1
            
            return dp(row + 1, column) + dp(row, column + 1)
        
        return dp(0, 0)