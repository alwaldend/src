class Solution:
    def findSpecialInteger(self, arr: List[int]) -> int:
        prev, count = arr[0], 1
        quarter = len(arr) / 4
        for num in arr[1:]:
            if num == prev:
                count += 1
            else:
                prev = num
                count = 1
            if count > quarter:
                break
        return prev