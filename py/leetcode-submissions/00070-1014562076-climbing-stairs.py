class Solution:
    def climbStairs(self, n: int) -> int:
        if n < 3:
            return n

        minus_one, minus_two = 2, 1
        for step in range(3, n + 1):
            minus_one, minus_two = minus_one + minus_two, minus_one
        
        return minus_one