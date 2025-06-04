class Solution:
    def findErrorNums(self, nums: List[int]) -> List[int]:
        length = len(nums)
        # dupl_xor_miss = duplicate ^ missing
        dupl_xor_miss = reduce(lambda total, i: total ^ i ^ nums[i - 1], range(length + 1))
        rightmost_set_bit = dupl_xor_miss & -dupl_xor_miss
        xor_group1 = xor_group2 = 0
        for i in range(1, length + 1):
            if i & rightmost_set_bit:
                xor_group1 ^= i
            else:
                xor_group2 ^= i
            if nums[i - 1] & rightmost_set_bit:
                xor_group1 ^= nums[i - 1]
            else:
                xor_group2 ^= nums[i - 1]
        for num in nums:
            if num == xor_group1:
                return num, xor_group2
            if num == xor_group2:
                return num, xor_group1 
        
        raise Exception()