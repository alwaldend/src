class Solution:
    def subsetXORSum(self, nums: List[int]) -> int:
        return reduce(operator.or_, nums) * 1 << (len(nums) - 1)
        