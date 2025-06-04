class Solution:
    # [0,3,2,1,0]
    def peakIndexInMountainArray(self, arr: List[int]) -> int:
        length = len(arr)

        left, right = 0, length - 1
        while left <= right:
            peak = left + (right - left) // 2
            left_val, right_val = arr[peak-1], arr[peak+1]
            peak_val = arr[peak]

            if left_val < peak_val > right_val:
                return peak
            
            if right_val > peak_val:
                left = peak + 1
            else:
                right = peak - 1
        
        return left
        