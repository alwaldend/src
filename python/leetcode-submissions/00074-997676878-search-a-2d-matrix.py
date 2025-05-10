class Solution:
    def searchMatrix(self, matrix: List[List[int]], target: int) -> bool:
        length_vertical, length_horizontal = len(matrix), len(matrix[0])


        left, right = 0, length_vertical - 1
        while left <= right:
            mid = left + (right - left) // 2
            mid_number = matrix[mid][0]
            
            if mid_number == target:
                return True
            
            if mid_number > target:
                right = mid - 1
            else:
                left = mid + 1
        
        vertical_index = right
        left, right = 0, length_horizontal - 1
        while left <= right:
            mid = left + (right - left) // 2
            mid_number = matrix[vertical_index][mid]

            if mid_number == target:
                return True
            
            if mid_number > target:
                right = mid - 1
            else:
                left = mid + 1
            
        return False