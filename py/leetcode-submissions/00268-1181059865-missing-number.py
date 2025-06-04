import operator
class Solution:
    def missingNumber(self, nums: List[int]) -> int:
        nums.extend(range(len(nums) + 1))
        return reduce(operator.xor, nums, 0)