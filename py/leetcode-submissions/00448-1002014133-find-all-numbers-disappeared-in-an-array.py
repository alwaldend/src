class Solution:
    # [1,1,3,4], [1,2,3,4] -> [2]
    def findDisappearedNumbers(self, nums: List[int]) -> List[int]:
        result = {number: True for number in range(1, len(nums) + 1)}

        for number in nums:
            result[number] = False
        
        return [number for number, valid in result.items() if valid]
