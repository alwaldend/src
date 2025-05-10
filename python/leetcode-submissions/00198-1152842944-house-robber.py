class Solution:
    def rob(self, nums: List[int]) -> int:
        length = len(nums)
        cache = {}

        def dp(i: int) -> int:
            if i >= length:
                return 0
            val = nums[i]
            if i in cache:
                return cache[i]
            res = max(dp(i + 1), val + dp(i + 2))
            cache[i] = res
            return res 

        return dp(0)