class Solution:
    def searchInsert(self, nums: List[int], target: int) -> int:
        length = len(nums)
        left, right = 0, length - 1

        while left <= right:
            mid = left + (right - left) // 2
            mid_number = nums[mid]

            if mid_number == target:
                return mid

            if mid_number > target:
                right = mid - 1
            else:
                left = mid + 1

        return left