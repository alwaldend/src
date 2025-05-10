class Solution:
    def canConstruct(self, ransomNote: str, magazine: str) -> bool:
        ransom_i, magazine_i = 0, 0
        ransom_length, magazine_length = len(ransomNote), len(magazine)

        if magazine_length < ransom_length:
            return False
        
        if ransom_length == 1:
            return ransomNote in magazine

        symbols = [0 for _ in range(26)]

        for symbol in magazine:
            symbols[ord(symbol)-97] += 1

        for symbol in ransomNote:
            symbols[ord(symbol)-97] -= 1
        
        
        return not any((count < 0 for count in symbols))