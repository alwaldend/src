class Solution:
    def findErrorNums(self, nums: List[int]) -> List[int]:
        n = len(nums)
        dupl_xor_miss = 0
        for i in range(1, n+1):
            dupl_xor_miss ^= i ^ nums[i-1]
        
        # example for get rightmost set bit
        # x:             01110000
        # ~x:            10001111
        # -x or ~x + 1:  10010000
        # x & -x:        00010000

        # example for unset rightmost set bit
        # x:             01110000
        # x-1:           01101111
        # x & (x-1):     01100000
        rightmost_set_bit = dupl_xor_miss & -dupl_xor_miss
        xor_group1 = xor_group2 = 0
        for i in range(1, n + 1):
            if i & rightmost_set_bit:
                xor_group1 ^= i
            else:
                xor_group2 ^= i
            if nums[i-1] & rightmost_set_bit:
                xor_group1 ^= nums[i-1]
            else:
                xor_group2 ^= nums[i-1]
        
        for num in nums:
            if num == xor_group1:
                return [num, xor_group2]
            if num == xor_group2:
                return [num, xor_group1]

        return []