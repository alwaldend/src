class Solution:
    def gameOfLife(self, board: List[List[int]]) -> None:
        """
        Do not return anything, modify board in-place instead.
        """
        rows_count, cols_count = len(board), len(board[0])
        moves = ((0, 1), (0, -1), (1, 0), (-1, 0), (1, 1), (1, -1), (-1, 1), (-1, -1))
        flip_cells = set()
        
        for row in range(rows_count):
            for col in range(cols_count):
                live_neighbors = 0
                is_alive = board[row][col] == 1

                for row_delta, col_delta in moves:
                    new_row, new_col = row + row_delta, col + col_delta
                    if not 0 <= new_row < rows_count or not 0 <= new_col < cols_count:
                        continue
                    if board[new_row][new_col] == 1:
                        live_neighbors += 1

                if (is_alive and (live_neighbors < 2 or live_neighbors > 3)) or (
                    not is_alive and live_neighbors == 3
                ):
                    flip_cells.add((row, col))
        
        while flip_cells:
            row, col = flip_cells.pop()
            board[row][col] = 0 if board[row][col] == 1 else 1