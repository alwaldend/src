class Solution:
    def minimumReplacement(self, nums: List[int]) -> int:
        nums_count = len(nums)
        operations_count = 0
        for i in reversed(range(nums_count - 1)):
            cur, prev = nums[i], nums[i+1]
            if cur <= prev:
                continue
            
            elements_count = (cur + prev - 1) // prev
            operations_count += elements_count - 1
            nums[i] //= elements_count

        return operations_count