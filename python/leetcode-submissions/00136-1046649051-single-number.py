class Solution:
    def singleNumber(self, nums: List[int]) -> int:
        return reduce(lambda total, element: total ^ element, nums)