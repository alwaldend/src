class Solution:
    def isHappy(self, n: int) -> bool:
        while True:
            sum = 0
            while n > 0:
                sum += (n % 10 )**2
                n = n // 10
            n = sum
            print(n)
            if sum < 10:
                break

        return n in [1, 7]