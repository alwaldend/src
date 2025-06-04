class Solution:
    def isMonotonic(self, nums: List[int]) -> bool:
        length = len(nums)
        if length < 3:
            return True
        is_increasing = None
        for i in range(1, length):
            cur, prev = nums[i], nums[i-1]
            if cur == prev:
                continue
            is_increasing_cur = cur > prev
            if is_increasing is None:
                is_increasing = is_increasing_cur
                continue
            if is_increasing and not is_increasing_cur or (
                not is_increasing and is_increasing_cur
            ):
                return False
        return True