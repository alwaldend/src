class Solution:
    def lengthOfLIS(self, nums: List[int]) -> int:
        nums_count = len(nums)
        dp = [0] * nums_count
        dp[0] = 1
        for i in range(1, nums_count):
            length = 0
            start = nums[i]

            for j in range(0, i):
                if start > nums[j]:
                    length = max(length, dp[j])

            dp[i] = 1 + length
            
        return max(dp)