class Solution:
    def firstPalindrome(self, words: List[str]) -> str:
        for s in words:
            for i in range(len(s) // 2):
                if s[i] != s[-i - 1]:
                    break
            else:
                return s
        return ""