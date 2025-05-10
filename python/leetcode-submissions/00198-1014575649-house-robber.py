class Solution:
    def rob(self, nums: List[int]) -> int:
        house_count = len(nums)

        prev_1, prev_2 = 0, 0

        for house in range(0, house_count):
            prev_1, prev_2 = max(prev_2 + nums[house], prev_1), prev_1

        return prev_1
