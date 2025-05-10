class Solution:
    def climbStairs(self, n: int) -> int:
        prev, cur = 1, 2
        if n == prev:
            return prev
        if n == cur:
            return cur
        for i in range(2, n):
            new = prev + cur
            cur, prev = new, cur
        return cur