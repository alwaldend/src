class Solution:
    def combinationSum4(self, nums: List[int], target: int) -> int:

        @cache
        def dp(cur_sum: int) -> int:
            if cur_sum == target:
                return 1
            if cur_sum > target:
                return 0
            
            return sum(dp(cur_sum + num) for num in nums)

        return dp(0)
            