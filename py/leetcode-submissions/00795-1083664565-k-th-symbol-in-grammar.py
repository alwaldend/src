class Solution:
    def kthGrammar(self, n: int, k: int) -> int:

        def get(row: int, column: int) -> int:
            if column == 0 or column == 1:
                return column
            
            prev_row_length = 2 ** (row - 1)
            if column >= prev_row_length:
                return 1 ^ get(row - 1, column - prev_row_length)
            return get(row - 1, column)

        return get(n - 1, k - 1)

        # 0
        # 0 1
        # 0 1 1 0
        # 0 1 1 0 1 0 0 1
        # 0 1 1 0 1 0 0 1
