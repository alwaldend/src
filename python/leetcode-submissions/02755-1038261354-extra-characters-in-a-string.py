class Solution:
    def minExtraChar(self, s: str, dictionary: List[str]) -> int:
        chars_count = len(s)

        @cache
        def dp(i: int) -> int:
            chars_left = chars_count - i
            if chars_left == 0:
                return 0
            
            min_extra_chars = chars_left

            for word in dictionary:
                word_length = len(word)

                if word_length > chars_left or s[i:i+word_length] != word:
                    continue
                
                if word_length == chars_left:
                    return 0

                min_extra_chars = min(min_extra_chars, dp(i + word_length))

            return min(min_extra_chars, 1 + dp(i + 1))
        
        return dp(0)

                
