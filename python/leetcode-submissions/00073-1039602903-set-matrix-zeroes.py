class Solution:
    def setZeroes(self, matrix: List[List[int]]) -> None:
        """
        Do not return anything, modify matrix in-place instead.
        """
        rows_count, cols_count = len(matrix), len(matrix[0])
        flip_rows, flip_cols = set(), set()

        for row in range(rows_count):
            for col in range(cols_count):
                if matrix[row][col] == 0:
                    flip_rows.add(row)
                    flip_cols.add(col)

        while flip_rows:
            row = flip_rows.pop()
            for col in range(cols_count):
                if matrix[row][col] != 0:
                    matrix[row][col] = 0
        
        while flip_cols:
            col = flip_cols.pop()
            for row in range(rows_count):
                if matrix[row][col] != 0:
                    matrix[row][col] = 0
        
            
