from sortedcontainers import SortedSet

class Solution:
    def longestConsecutive(self, nums: List[int]) -> int:
        nums = SortedSet(nums)
        longest, current = 0, 0
        for number in nums:
            if number - 1 in nums:
                current += 1
                continue
            if current > longest:
                longest = current
            current = 1
        return max(current, longest)
            