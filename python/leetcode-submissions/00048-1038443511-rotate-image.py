class Solution:
    def rotate(self, matrix: List[List[int]]) -> None:
        """
        Do not return anything, modify matrix in-place instead.
        """
        side_length = len(matrix)

        for row in range(side_length // 2):
            matrix[row], matrix[-row-1] = matrix[-row-1], matrix[row]
        
        for i in range(side_length):
            for j in range(i):
                matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]