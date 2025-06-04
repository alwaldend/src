class Solution:
    def minimumDeleteSum(self, s1: str, s2: str) -> int:
        length_1, length_2 = len(s1), len(s2)

        @cache
        def calculate(i: int, j: int) -> int:    
            if i >= length_1 and j >= length_2:
                return 0
            
            if i >= length_1:
                return sum(ord(char) for char in s2[j:])
            
            if j >= length_2:
                return sum(ord(char) for char in s1[i:])
            
            if s1[i] == s2[j]:
                return calculate(i + 1, j + 1)
            
            return min(
                ord(s1[i]) + calculate(i + 1, j),
                ord(s2[j]) + calculate(i, j + 1)
            )

        return calculate(0, 0)