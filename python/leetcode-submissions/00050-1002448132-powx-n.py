class Solution:
    def myPow(self, x: float, n: int) -> float:
        if n == 0:
            return 1

        if n < 0:
            n *= -1
            x = 1 / x

        result = 1
        while n:
            if n % 2:
                result *= x
                n -= 1
            x *= x
            n //= 2
        
        return result