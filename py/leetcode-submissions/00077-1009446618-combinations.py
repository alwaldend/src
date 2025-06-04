class Solution:
    def combine(self, n: int, k: int) -> List[List[int]]:
        current, result = [], []
        
        def backtrack(first: int) -> None:
            if len(current) == k:
                result.append(tuple(current[:]))
                return

            for i in range(first, n + 1):
                current.append(i)
                backtrack(i + 1)
                current.pop()
            
            return

        backtrack(1)

        return result 