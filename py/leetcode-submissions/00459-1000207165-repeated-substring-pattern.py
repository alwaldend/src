class Solution:
    def repeatedSubstringPattern(self, s: str) -> bool:
        length = len(s)
        for i in range(length-1):
            if length % (i + 1) != 0:
                continue
            
            if s == s[:i+1] * (length // (i + 1)):
                return True
        
        return False