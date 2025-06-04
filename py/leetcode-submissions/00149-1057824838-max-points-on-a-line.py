class Solution:
    def maxPoints(self, points: List[List[int]]) -> int:
        lines = defaultdict(set)
        length = len(points)
        if length < 3:
            return length
        for i in range(length):
            x1, y1 = points[i]
            for j in range(i + 1, length):
                x2, y2 = points[j]
                k1, k2 = None, None
                if y1 == y2:
                    k2 = y1
                elif x1 == x2:
                    k1 = x1
                else: 
                    k1 = (y2 - y1) / (x2 - x1)
                    k2 = y1 - k1 * x1
                lines[(k1, k2)].update(((x1, y1), (x2, y2)))
        return max((len(points) for points in lines.values()), default=0)
