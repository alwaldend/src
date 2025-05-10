class Solution:
    def singleNumber(self, nums: List[int]) -> int:
        nums.sort()
        length = len(nums)
        for i in range(0, length, 3):
            if i == length - 1:
                return nums[i]
    
            num_1, num_2, num_3 = nums[i], nums[i+1], nums[i+2]

            if num_1 == num_2 == num_3:
                continue
            
            if num_2 == num_3:
                return num_1
            
            if num_1 == num_3:
                return num_2
            
            return num_3
            