class Solution:
    def updateMatrix(self, mat: List[List[int]]) -> List[List[int]]:
        if not mat or not mat[0]:
            return []

        row_count, column_count = len(mat), len(mat[0])
        queue = set()
        MAX_VALUE = row_count * column_count
        
        # Initialize the queue with all 0s and set cells with 1s to MAX_VALUE.
        for row in range(row_count):
            for column in range(column_count):
                if mat[row][column] == 0:
                    queue.add((row, column))
                else:
                    mat[row][column] = MAX_VALUE
        
        directions = ((1, 0), (-1, 0), (0, 1), (0, -1))
        
        while queue:
            center_row, center_column = queue.pop()
            distance_from_center = mat[center_row][center_column] + 1

            for delta_row, delta_column in directions:
                new_row, new_column = center_row + delta_row, center_column + delta_column
                if not 0 <= new_row < row_count or not 0 <= new_column < column_count:
                    continue
    
                if mat[new_row][new_column] > distance_from_center:
                    queue.add((new_row, new_column))
                    mat[new_row][new_column] = distance_from_center
        
        return mat