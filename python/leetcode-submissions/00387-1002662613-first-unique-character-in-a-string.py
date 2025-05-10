class Solution:
    def firstUniqChar(self, s: str) -> int:
        counts = {}
        repeated_index = len(s)
        for i, letter in enumerate(s):
            if letter in counts:
                counts[letter] = repeated_index
                continue
            
            counts[letter] = i
        
        result = repeated_index

        for letter, index in counts.items():
            if index == repeated_index:
                continue
            
            if index < result:
                result = index

        return result if result != repeated_index else -1