class Solution:
    def numDecodings(self, s: str) -> int:
        char_count = len(s)

        @cache
        def dfs(i: int) -> int:
            if i == char_count: 
                return 1
            if s[i] == "0": 
                return 0
            return dfs(i + 1) + (
                dfs(i + 2) if i + 1 < char_count and s[i:i+2] < "27" else 0
            )

        return dfs(0)