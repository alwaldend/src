class Solution:
    def lengthOfLongestSubstring(self, s: str) -> int:
        length = len(s)
        if length < 2:
            return length
        
        max_length = 0
        start = 0
        indexes = {s[0]: 0}
        for end in range(1, length):
            letter = s[end]
            if letter not in indexes:
                indexes[letter] = end
                continue

            this_length = end - start
            if this_length > max_length:
                max_length = this_length
            
            letter_index = indexes[letter]
            
            for remove_letter in s[start:letter_index]:
                indexes.pop(remove_letter)

            indexes[letter], start = end, letter_index + 1

        return max(max_length, length - start)