class Solution:
    def minOperations(self, nums: List[int], k: int) -> int:
        if not nums:
            return -1
        nums.sort(reverse=True)
        if nums[-1] < k:
            return -1
        prev = nums[0]
        if prev == k:
            return 0
        ops = 0
        for num in nums[1:]:
            if num == prev:
                continue
            ops += 1
            prev = num
            if num == k:
                return ops
        return ops + 1