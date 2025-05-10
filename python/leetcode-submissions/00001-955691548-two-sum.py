class Solution:
    def twoSum(self, nums: List[int], target: int) -> List[int]:
        for i, number in enumerate(nums):
            for j, number_inner in enumerate(nums[i+1:]):
                if number + number_inner == target:
                    return [i, i+j+1]