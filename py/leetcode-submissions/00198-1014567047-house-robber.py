class Solution:
    def rob(self, nums: List[int]) -> int:
        house_count = len(nums)
        
        @cache
        def dp(house: int) -> int:
            if house >= house_count:
                return 0

            return max(dp(house + 1), nums[house] + dp(house + 2)) 

        return dp(0)
