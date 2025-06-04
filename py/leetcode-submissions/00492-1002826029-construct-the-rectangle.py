class Solution:
    def constructRectangle(self, area: int) -> List[int]:
        return next([area//width, width] 
                    for width in range(int(area**0.5), 0, -1) 
                    if area % width == 0)