import operator
class Solution:
    def missingNumber(self, nums: List[int]) -> int:
        return reduce(
            operator.xor, 
            nums, 
            reduce(operator.xor, tuple(range(len(nums) + 1)), 0)
        )