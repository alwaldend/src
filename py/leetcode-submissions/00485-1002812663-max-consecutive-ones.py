class Solution:
    def findMaxConsecutiveOnes(self, nums: List[int]) -> int:
        max_count, count = 0, 0
        for number in nums:
            if not number:
                max_count = max(max_count, count)
                count = 0
                continue
            
            count += 1
        
        return max(max_count, count)