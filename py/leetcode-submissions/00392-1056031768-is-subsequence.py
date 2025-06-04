class Solution:
    def isSubsequence(self, s: str, t: str) -> bool:
        i = 0
        length = len(s)
        for char in t:
            if i == length:
                break
            if s[i] == char:
                i += 1
            
        return i == length