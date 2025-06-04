class Solution:
    # two indexes: ransom_i, magazine_i
    # sort both ransomNote and magazine (ascending order)
    # while indexes have not reached the end:
    # - if ransom symbol is equal to the magazine_symbol: increase both indexes
    # - if it is not equal, increase only magazine index
    # return wheter the ranson index is equal to the length of the ransom
    def canConstruct(self, ransomNote: str, magazine: str) -> bool:
        ransom_i, magazine_i = 0, 0
        ransom_length, magazine_length = len(ransomNote), len(magazine)

        if magazine_length < ransom_length:
            return False
        
        if ransom_length == 1:
            return ransomNote in magazine

        ransomNote = sorted(ransomNote)
        magazine = sorted(magazine)

        while ransom_i < ransom_length and magazine_i < magazine_length:
            ransom_symbol, magazine_symbol = ransomNote[ransom_i], magazine[magazine_i]
            
            if magazine_symbol > ransom_symbol:
                return False

            if ransom_symbol == magazine_symbol:
                ransom_i += 1
            
            magazine_i += 1

        return ransom_i == ransom_length