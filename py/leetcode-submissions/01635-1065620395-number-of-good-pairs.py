class Solution:
    def numIdenticalPairs(self, nums: List[int]) -> int:
        counter = Counter(nums)
        
        return sum((count ** 2 - count) // 2 for count in counter.values() if count > 1)