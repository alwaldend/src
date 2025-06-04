class Solution:
    def firstUniqChar(self, s: str) -> int:
        counts = defaultdict(int)
        repeated_index = len(s)
        for letter in s:
            counts[letter] += 1
        
        for i, letter in enumerate(s):
            if counts[letter] != 1:
                continue
            
            return i
        
        return -1