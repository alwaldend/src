class Solution:
    def canMakeSubsequence(self, str1: str, str2: str) -> bool:
        length1, length2 = len(str1), len(str2)
        
        for i in range(length1):
            if i + length2 > length1:
                break
            
            idx2 = 0
            for idx1 in range(i, length1):
                if ord(str2[idx2]) - ord(str1[idx1]) in (0, 1, -25):
                    idx2 += 1
                
                if idx2 == length2:
                    return True
            
            
        return False