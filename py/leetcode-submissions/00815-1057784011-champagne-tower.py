class Solution:
    def champagneTower(self, poured: int, query_row: int, query_glass: int) -> float:
        if query_row == 0:
            return min(1, poured)
        if query_glass in (0, query_row):
            while query_row and poured:
                poured = (poured - 1) / 2 if poured > 1 else 0
                query_row -= 1
            return min(1, poured)

        cur_row, next_row = [0] * (query_row + 2), [0] * (query_row + 2)
        cur_row[0] = poured
        for row in range(query_row + 1):
            next_row[0] = 0
            for col in range(row + 1):
                overflow = (cur_row[col] - 1) / 2
                if overflow > 0:
                    next_row[col] += overflow
                    next_row[col + 1] = overflow
                else:
                    next_row[col + 1] = 0
            cur_row, next_row = next_row, cur_row
        return min(1, next_row[query_glass])