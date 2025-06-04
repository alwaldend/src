class Solution:
    def search(self, nums: List[int], target: int) -> int:
        nums_count = len(nums)
        left, right = 0, len(nums) - 1

        first_num, last_num = nums[0], nums[-1]
        if first_num == target:
            return 0
        if last_num == target:
            return nums_count - 1
        
        # Find the index of the pivot element (the smallest element)
        while left <= right:
            mid = left + (right - left) // 2
            if nums[mid] > last_num:
                left = mid + 1
            else:
                right = mid - 1
        
        pivot_num = nums[left]
        if pivot_num == target:
            return left
        
        if pivot_num < target < last_num:
            right = nums_count - 1
        else:
            left = 0
        
        while left <= right:
            mid = left + (right - left) // 2
            mid_num = nums[mid]
            if mid_num == target:
                return mid
            elif mid_num > target:
                right = mid - 1
            else:
                left = mid + 1
        
        return -1