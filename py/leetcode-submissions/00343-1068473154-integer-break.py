class Solution:
    def integerBreak(self, n: int) -> int:
        if n < 4:
            return n - 1
        
        @cache
        def dp(num: int) -> int:
            if num <= 3:
                return num
            ans = num
            for i in range(2, num):
                ans = max(ans, i * dp(num - i))
            return ans

        return dp(n)