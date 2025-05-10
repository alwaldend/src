class Solution:
    def canJump(self, nums: List[int]) -> bool:
        length = len(nums)

        if length < 2:
            return True
        
        current = nums[0]

        for i in range(1, length):
            if current == 0:
                return False
            current -= 1
            current = max(current, nums[i])
        
        return True
            