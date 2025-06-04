class Solution:
    # y x x
    # x x x
    # x x x: (2 / 8) * (4 / 16) = 0.0625
    #        (1/8, 1/8), (2/8, 2/8)
    # 
    # x x x
    # x y x
    # x x x: 0
    def __init__(self):
        self.available_moves = (
            (2, 1), (1, 2), (-2, 1), (-1, 2), (2, -1), (1, -2), 
            (-2, -1), (-1, -2)
        )
    

    def knightProbability(self, n: int, k: int, row: int, column: int) -> float:
        if k < 1:
            return 1
        
        if n < 3:
            return 0

        return self.calculate(row, column, n, k)
    
    @cache
    def calculate(self, row: int, column: int, size: int, 
                       moves_left: int) -> float:
       
        if moves_left < 1:
            return 1

        probability = 0

        for row_add, column_add in self.available_moves:
            new_row = row_add + row
            new_column = column_add + column

            if new_row < 0 or new_row >= size or (
                new_column < 0 or new_column >= size
            ):
                continue

            probability += self.calculate(
                new_row, new_column, size, moves_left - 1
            ) / 8

        return probability