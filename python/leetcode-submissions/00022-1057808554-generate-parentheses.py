class Solution:
    def generateParenthesis(self, n: int) -> List[str]:

        current = []
        current_max = n * 2
        chars = "()"


        def backtrack(open: int, closed: int) -> Generator[None, None, List[str]]:
            if len(current) == current_max:
                yield "".join(current)
            
            if open:
                current.append(chars[0])
                yield from backtrack(open - 1, closed)
                current.pop()
            
            if closed and closed > open:
                current.append(chars[1])
                yield from backtrack(open, closed - 1)
                current.pop()
        
        return tuple(combination for combination in backtrack(n, n))