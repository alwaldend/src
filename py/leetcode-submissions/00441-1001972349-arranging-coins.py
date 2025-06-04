class Solution:
    def arrangeCoins(self, n: int) -> int:
        count = 0
        row = 1
        while n >= row:
            n -= row
            row += 1
            count += 1

        return count