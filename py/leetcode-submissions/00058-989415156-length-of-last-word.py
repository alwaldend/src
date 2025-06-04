class Solution:
    # we need to iterate over the string searching for words
    # possible combinations:
    # 1. a whitespace after a symbol: end of the word
    # 2. a whitespace after a whitespace: ignore
    # 3. a symbol after a symbol, but it's the last symbol: end of the last word
    # 4. a symbol after a symbol: ignore
    def lengthOfLastWord(self, s: str) -> int:
        length = len(s)
        if length == 1:
            return 1

        word_start, word_end, in_word = 0, 0, False

        for i, symbol in enumerate(s):
            is_whitespace = symbol == " "
            is_last = i == length - 1
            
            if is_whitespace and in_word:
                in_word = False
                word_end = i - 1
                continue
            
            if is_whitespace and not in_word:
                continue

            if not is_whitespace and in_word and is_last:
                word_end = i
                continue

            if not is_whitespace and not in_word and is_last:
                word_end = i
                word_start = i
                continue

            if not is_whitespace and in_word:
                continue

            if not is_whitespace and not in_word:
                word_start = i
                in_word = True
                continue
            
        return word_end - word_start + 1