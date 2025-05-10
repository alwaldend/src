class Solution:
    def spiralOrder(self, matrix: List[List[int]]) -> List[int]:
        row_count, col_count = len(matrix), len(matrix[0])
        if row_count == 1:
            return matrix[0]
        if col_count == 1:
            return [row[0] for row in matrix]


        result = []
        top, bot, left, right = 0, row_count - 1, 0, col_count - 1
        capacity = row_count * col_count

        while len(result) < capacity:
            for col in range(left, right + 1):
                result.append(matrix[top][col])
            top += 1

            for row in range(top, bot + 1):
                result.append(matrix[row][right])
            right -= 1

            if top <= bot:
                for col in range(right, left - 1, -1):
                    result.append(matrix[bot][col])
                bot -= 1
            
            if left <= right:
                for row in range(bot, top - 1, -1):
                    result.append(matrix[row][left])
                left += 1

        return result