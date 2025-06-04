class Solution:
    def trailingZeroes(self, n: int) -> int:
        zeros = 0
        while n != 0:
            new_n = n // 5
            zeros += new_n
            n = new_n
        return zeros
        
