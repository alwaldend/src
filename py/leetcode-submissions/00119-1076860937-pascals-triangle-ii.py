class Solution:
    def getRow(self, rowIndex: int) -> List[int]:
        cur, prev = [], [1]
        row = 0

        while row < rowIndex:
            cur.append(1)
            for i in range(1, len(prev)):
                cur.append(prev[i] + prev[i-1])
            cur.append(1)
            prev.clear()
            cur, prev = prev, cur
            row += 1
        
        return prev 