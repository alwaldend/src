class DetectSquares:

    def __init__(self):
        self._row_col = defaultdict(lambda: defaultdict(int))

    def add(self, point: List[int]) -> None:
        self._row_col[point[0]][point[1]] += 1

    def count(self, point: List[int]) -> int:
        ways_count = 0
        row1, col1 = point
        for col2, col2_count in self._row_col[row1].items():
            if col2 == col1:
                continue

            side = col2 - col1
            for row2 in (row1 + side, row1 - side):
                point3_count = self._row_col[row2][col1]
                point4_count = self._row_col[row2][col2]
                ways_count += col2_count * point3_count * point4_count

        return ways_count

# Your DetectSquares object will be instantiated and called as such:
# obj = DetectSquares()
# obj.add(point)
# param_2 = obj.count(point)