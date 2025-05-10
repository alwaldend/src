class Solution:
    def productExceptSelf(self, nums: List[int]) -> List[int]:
        nums_count = len(nums)
        result = [1] * nums_count
        prefix = 1
        postfix = 1
        for i in range(nums_count):
            result[i] *= prefix
            prefix *= nums[i]
            result[nums_count-i-1] *= postfix
            postfix *= nums[nums_count-i-1]
        
        return result