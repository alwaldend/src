class Solution:
    def minOperations(self, nums: List[int], x: int) -> int:
        target, length = sum(nums) - x, len(nums)
        max_len = cur_sum = left = 0
        
        if target == 0:
            return length
        
        for right, val in enumerate(nums):
            cur_sum += val
            while left <= right and cur_sum > target:
                cur_sum -= nums[left]
                left += 1
            if cur_sum == target:
                max_len = max(max_len, right - left + 1)
        
        return length - max_len if max_len else -1