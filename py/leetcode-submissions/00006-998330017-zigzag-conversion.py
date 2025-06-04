class Solution:
    def convert(self, s: str, numRows: int) -> str:
        if numRows < 2:
            return s
        
        result: List[List[str]] = [[] for _ in range(numRows)]
    
        row, is_ascending, last = 0, True, numRows - 1
        for symbol in s:
            result[row].append(symbol)
            
            if is_ascending and row < last:
                row += 1
            elif is_ascending:
                is_ascending = False
            
            if is_ascending:
                continue

            if row > 0:
                row -= 1
            else:
                is_ascending = True
                row = 1
    
        
        return "".join("".join(row) for row in result)