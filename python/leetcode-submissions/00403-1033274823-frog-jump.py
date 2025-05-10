class Solution:
    def canCross(self, stones: List[int]) -> bool:
        m = set(stones)
        @cache
        def dfs(i, j):
            if i == stones[-1]: return True
            return any(x and x + i in m and dfs(x + i, x) for x in range(j - 1, j + 2))
        return dfs(0, 0)