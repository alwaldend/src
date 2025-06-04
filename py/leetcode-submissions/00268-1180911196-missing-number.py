class Solution:
    def missingNumber(self, nums: List[int]) -> int:
        def xor(total: int, i: int) -> int:
            return total ^ nums[i] ^ (i + 1)
        return reduce(xor, chain((0, ), range(len(nums))))