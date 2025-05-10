class Solution:
    def repeatedSubstringPattern(self, s: str) -> bool:
        length = len(s)
        for i in range(1, length // 2 + 1):
            if length % i != 0:
                continue
            
            if s == s[:i] * (length // i):
                return True
        
        return False