class Solution:
    def climbStairs(self, n: int) -> int:

        @cache
        def dp(step: int) -> int:
            if step <= 2:
                return step
            
            return dp(step - 1) + dp(step - 2) 
 
        return dp(n)