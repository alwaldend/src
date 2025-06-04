class Solution:
    def isPowerOfTwo(self, n: int) -> bool:
        if n == 1:
            return True

        while n > 2 and not n % 2:
            n //= 2
        return n == 2