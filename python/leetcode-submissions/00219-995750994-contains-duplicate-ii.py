class Solution:
    def containsNearbyDuplicate(self, nums: List[int], k: int) -> bool:
        indexes = {}
        for i, number in enumerate(nums):
            if number in indexes and abs(i - indexes[number]) <= k:
                return True
            indexes[number] = i
        return False

            