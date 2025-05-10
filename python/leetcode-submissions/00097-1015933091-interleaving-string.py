class Solution:
    def isInterleave(self, s1: str, s2: str, s3: str) -> bool:
        length_1, length_2, length_3 = len(s1), len(s2), len(s3)
        if length_1 + length_2 != length_3:
            return False
        
        @cache
        def dp(i_1: int, i_2: int, i_3: int) -> bool:
            if i_3 == length_3:
                return True

            target = s3[i_3]
            return (
                i_1 != length_1 and s1[i_1] == target and dp(i_1 + 1, i_2, i_3 + 1)
            ) or (
                i_2 != length_2 and s2[i_2] == target and dp(i_1, i_2 + 1, i_3 + 1)
            )

        return dp(0, 0, 0)
