class Solution:
    def singleNumber(self, nums: List[int]) -> int:
        nums_count = len(nums)
        if nums_count == 1:
            return nums[0]
        
        nums.sort()
        for i in range(0, nums_count, 2):
            if i + 1 == nums_count:
                return nums[i]

            cur_num, next_num = nums[i], nums[i+1]
            if cur_num == next_num:
                continue
            
            return cur_num
        
        return 0