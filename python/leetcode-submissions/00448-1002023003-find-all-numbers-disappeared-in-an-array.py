class Solution:
    # [1,1,3,4], [1,2,3,4] -> [2]
    def findDisappearedNumbers(self, nums: List[int]) -> List[int]:
        for number in nums:
            index = abs(number) - 1
            nums[index] = -1 * abs(nums[index])
        
        return [number
                for number in range(1, len(nums) + 1) 
                if nums[number-1] > 0]
        
