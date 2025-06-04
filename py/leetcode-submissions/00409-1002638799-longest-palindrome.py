class Solution:
    def longestPalindrome(self, s: str) -> int:
        counts = defaultdict(int)

        for letter in s:
            counts[letter] += 1
        
        result = 0
        used_odd_letter = False
        for letter, count in counts.items():
            is_odd = count % 2 != 0 
            if is_odd and used_odd_letter:
                result -= 1
            elif is_odd:
                used_odd_letter = True
            result += count

        return result