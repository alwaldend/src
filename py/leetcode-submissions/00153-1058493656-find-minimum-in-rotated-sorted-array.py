class Solution:
    def findMin(self, nums: List[int]) -> int:
        length = len(nums)
        if length == 1:
            return nums[0]
        first, last = nums[0], nums[-1]
        if first < last:
            return first
        left, right = 0, length - 1
        while left <= right:
            mid = left + (right - left) // 2
            left_num = nums[mid - 1] if mid > 0 else first
            mid_num = nums[mid] 
            right_num = nums[mid + 1] if mid + 1 < length else last
            if left_num > mid_num <= right_num:
                return mid_num
            if left_num > mid_num or (left_num <= mid_num and mid_num >= first):
                left = mid + 1
            else:
                right = mid
        return nums[left]

