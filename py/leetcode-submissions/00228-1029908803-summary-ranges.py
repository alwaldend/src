class Solution:
    def summaryRanges(self, nums: List[int]) -> List[str]:
        nums_count = len(nums)
        if nums_count == 0:
            return []
        if nums_count == 1:
            return [str(nums[0])]

        ranges = [[nums[0]] * 2]
        for i, num in enumerate(nums[1:]):
            if ranges[-1][1] == num - 1:
                ranges[-1][1] = num
            else:
                ranges.append([num, num])

        return [f"{start}->{end}" if start != end else str(start) for start, end in ranges]