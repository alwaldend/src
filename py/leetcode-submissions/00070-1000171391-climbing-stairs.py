class Solution:
    # 1: 1
    # 2: 2
    # 3: 3 [2(2), 1(1)]
    # 4: 5 [3(3), (2)]
    def climbStairs(self, n: int) -> int:
        if n < 3:
            return n

        count, count_prev = 2, 1
        for number in range(3, n + 1):
            count, count_prev = count + count_prev, count
        return count