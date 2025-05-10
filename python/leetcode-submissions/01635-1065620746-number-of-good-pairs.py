class Solution:
    def numIdenticalPairs(self, nums: List[int]) -> int:        
        return sum((count ** 2 - count) // 2 for count in Counter(nums).values() if count > 1)