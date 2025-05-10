class Solution:
    def minFallingPathSum(self, matrix: List[List[int]]) -> int:
        length = len(matrix)
        max_val = 101 * length * length
        cache = {}
        deltas = ((1, 0), (1, 1), (1, -1))
        last_row = length - 1

        def dp(row: int, col: int) -> int:
            if (row, col) in cache:
                return cache[(row, col)]
            val = matrix[row][col]
            if row == last_row:
                return val
            min_cost = max_val
            for delta_row, delta_col in deltas:
                new_row, new_col = row + delta_row, col + delta_col
                if new_col == length or new_col == -1:
                    continue
                cost = dp(new_row, new_col)
                min_cost = min(min_cost, cost)
            res = val + min_cost
            cache[(row, col)] = res
            return res

        return min(dp(0, col) for col in range(length))