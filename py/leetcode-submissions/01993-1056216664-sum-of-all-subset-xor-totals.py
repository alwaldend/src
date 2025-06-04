class Solution:
    def subsetXORSum(self, nums: List[int]) -> int:
        subsets_count = 2**(len(nums) - 1)
        return reduce(operator.or_, nums) * subsets_count
        