import operator
class Solution:
    def missingNumber(self, nums: List[int]) -> int:
        length = len(nums)
        return reduce(operator.xor, range(length), reduce(operator.xor, nums, length))