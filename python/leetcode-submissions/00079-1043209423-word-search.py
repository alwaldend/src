class Solution:
    def exist(self, board: List[List[str]], word: str) -> bool:
        row_count, col_count = len(board), len(board[0])
        word_length = len(word)
        delta = ((0, 1), (0, -1), (1, 0), (-1, 0))

        if word_length == 1:
            return any(word in row for row in board)

        def backtrack(row: int, col: int, visited: Set, target: int) -> bool:
            if (row, col) in visited or not 0 <= row < row_count or not 0 <= col < col_count:
                return False

            if board[row][col] != word[target]:
                return False
            
            if target == word_length - 1:
                return True

            visited.add((row, col))
            if any(backtrack(row + delta_row, col + delta_col, visited, target + 1)
                   for delta_row, delta_col in delta):
                return True
            visited.remove((row, col))
            return False
        
        start_char = word[0]

        return any(backtrack(row, col, set(), 0) 
                   for row in range(row_count) 
                   for col in range(col_count)
                   if board[row][col] == start_char)