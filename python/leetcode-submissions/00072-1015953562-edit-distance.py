class Solution:
    def minDistance(self, word1: str, word2: str) -> int:
        length1, length2 = len(word1), len(word2)

        @cache
        def dp(i1: int, i2: int) -> int:
            if i1 == length1:
                return length2 - (i2 + 1)

            if i2 == length2:
                return length1 - (i1 + 1)
            
            if word1[i1] == word2[i2]:
                return dp(i1 + 1, i2 + 1)

            return 1 + min((
                # replace or insert
                dp(i1 + 1, i2 + 1),
                # remove from i1
                dp(i1 + 1, i2),
                # remove from i2
                dp(i1, i2 + 1)
            ))
        
        return dp(0, 0) + 1