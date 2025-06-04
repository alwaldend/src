class NumArray:

    def __init__(self, nums: List[int]):
        self._sums = tuple(accumulate(nums))

    def sumRange(self, left: int, right: int) -> int:
        # [-2, 0, 3, -5, 2, -1], [-2, -2, 1, -4, -2, -3]
        # [0, 2] -> 1
        if left == 0:
            return self._sums[right]

        return self._sums[right] - self._sums[left-1]


# Your NumArray object will be instantiated and called as such:
# obj = NumArray(nums)
# param_1 = obj.sumRange(left,right)