class Solution:
    def uniquePathsWithObstacles(self, obstacleGrid: List[List[int]]) -> int:
        row_count = len(obstacleGrid)
        column_count = len(obstacleGrid[0])
        
        if obstacleGrid[0][0] == 1 or obstacleGrid[-1][-1] == 1:
            return 0

        dp = [[0] * column_count for _ in range(row_count)]
        dp[0][0] = 1

        for column in range(1, column_count):
            if obstacleGrid[0][column] == 1:
                break

            dp[0][column] = dp[0][column-1]

        for row in range(1, row_count):
            if obstacleGrid[row][0] == 0:
                dp[row][0] = dp[row-1][0]

            for column in range(1, column_count):
                if obstacleGrid[row][column] == 1:
                    continue
                
                dp[row][column] = dp[row-1][column] + dp[row][column-1]
                     
        return dp[-1][-1]