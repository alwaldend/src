class Solution:
    def exist(self, board: List[List[str]], word: str) -> bool:
        row_count, col_count = len(board), len(board[0])
        word_length = len(word)
        delta = ((0, 1), (0, -1), (1, 0), (-1, 0))

        if word_length == 1:
            return any(word in row for row in board)
        if word_length > row_count * col_count:
            return False

        visited = set()

        def backtrack(row: int, col: int, target: int) -> bool:
            if target == word_length:
                return True
            if not 0 <= row < row_count or not 0 <= col < col_count or (
                (row, col) in visited or board[row][col] != word[target]
            ):
                return False

            visited.add((row, col))
            answer = any(backtrack(row + delta_row, col + delta_col, target + 1)
                         for delta_row, delta_col in delta)
            visited.remove((row, col))
            return answer

        return any(backtrack(row, col, 0) 
                   for row in range(row_count) 
                   for col in range(col_count))