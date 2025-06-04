class Solution:
    def sortByBits(self, arr: List[int]) -> List[int]:
        return tuple(num for num in sorted(arr, key=lambda num: (num.bit_count(), num)))