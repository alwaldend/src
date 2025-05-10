class Solution:
    def lengthOfLongestSubstring(self, s: str) -> int:
        length = len(s)
        if length < 2:
            return length
        
        max_length, left, charset = 1, 0, set([s[0]])
        for right in range(1, length):
            letter = s[right]
            if letter not in charset:
                charset.add(letter)
                continue

            this_length = right - left
            if this_length > max_length:
                max_length = this_length
            
            while letter in charset:
                charset.remove(s[left])
                left += 1
            
            charset.add(letter)

        return max(max_length, length - left)