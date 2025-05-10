class Solution:
    def maximalSquare(self, matrix: List[List[str]]) -> int:
        row_count, column_count = len(matrix), len(matrix[0])
        
        dp = [[0] * column_count for _ in range(row_count)]
        
        max_size = 0
        for column in range(column_count):
            if matrix[0][column] == "0":
                continue
            dp[0][column] = 1
            max_size = 1

        for row in range(row_count):
            if matrix[row][0] == "0":
                continue
            dp[row][0] = 1
            max_size = 1
        
        for row in range(1, row_count):
            for column in range(1, column_count):
                if matrix[row][column] == "0":
                    continue
                
                value = min(dp[row-1][column], dp[row][column-1], dp[row-1][column-1]) + 1
                dp[row][column] = value
                max_size = max(max_size, value)
        
        return max_size * max_size
