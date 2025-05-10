class Solution:
    def solve(self, board: List[List[str]]) -> None:
        """
        Do not return anything, modify board in-place instead.
        """
        row_count, col_count = len(board), len(board[0])
        moves = ((0, 1), (0, -1), (1, 0), (-1, 0))
        ignore_stack = set()

        def bfs(row: int, col: int) -> None:
            if not 0 <= row < row_count or not 0 <= col < col_count:
                return
            if (row, col) in ignore_stack or board[row][col] == "X":
                return

            ignore_stack.add((row, col))
            for row_delta, col_delta in moves:
                bfs(row + row_delta, col + col_delta)
        
        for row, col in itertools.chain(
            ((0, col) for col in range(col_count)),
            ((row_count - 1, col) for col in range(col_count)),
            ((row, 0) for row in range(row_count)),
            ((row, col_count - 1) for row in range(row_count))
        ):
            bfs(row, col)
        
        flip_x = True
        for row in range(1, row_count - 1):
            for col in range(1, col_count - 1):
                if (row, col) in ignore_stack:
                    continue

                board[row][col] = "X"
        