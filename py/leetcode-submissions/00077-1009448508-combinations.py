class Solution:
    def combine(self, n: int, k: int) -> List[List[int]]:
        current = []
        
        def backtrack(first: int) -> Generator[None, None, List[int]]:
            if len(current) == k:
                yield tuple(current[:])
                return

            for i in range(first, n + 1):
                current.append(i)
                yield from backtrack(i + 1)
                current.pop()
            
            return

        return tuple(combination for combination in backtrack(1)) 