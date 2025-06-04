class Solution:
    def containsNearbyDuplicate(self, nums: List[int], k: int) -> bool:
        indexes = defaultdict(set)
        for i, number in enumerate(nums):
            if number in indexes and any(abs(i - j) <= k for j in indexes[number]):
                return True
            indexes[number].add(i)
        return False

            