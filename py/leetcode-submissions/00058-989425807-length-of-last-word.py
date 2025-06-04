class Solution:
    # we need to iterate over the string searching for words
    # add a whitespace to the end to avoid the situation when the last symbol is non-whitespace 
    # possible combinations:
    # 1. a whitespace after a symbol: end of the word
    # 2. a whitespace after a whitespace: ignore
    # 3. a symbol after a symbol: ignore
    def lengthOfLastWord(self, s: str) -> int:
        s += " "
        result, word_length = 0, 0

        for i, symbol in enumerate(s):
            if symbol != " ":
                word_length += 1
                continue
            
            if word_length > 0:
                result = word_length
                word_length = 0
        
        return result