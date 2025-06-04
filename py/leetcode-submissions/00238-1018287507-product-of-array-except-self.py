class Solution:
    def productExceptSelf(self, nums: List[int]) -> List[int]:
        nums_count = len(nums)
        result = [1] * nums_count
        prefix = 1
        postfix = 1
        for i in range(nums_count):
            result[i] *= prefix
            prefix *= nums[i]

            from_end = -1 * (i + 1)
            result[from_end] *= postfix
            postfix *= nums[from_end]
        
        return result