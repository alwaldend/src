class Solution:
    def isPowerOfFour(self, n: int) -> bool:
        return n and not (n & (n - 1)) and (n & 0x55555555)