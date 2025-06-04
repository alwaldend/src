class Solution:
    # check edge cases:
    # - if length of the needle is less than the length of the haystack: -1
    # - if strings are equal: 0
    # - if length of the needle is equal to the length of the haystack, but 
    #   strings are not equal: -1
    # iterate over haystack and needle cheking if a symbol from haystack 
    #   corresponds to the symbol in needle
    def strStr(self, haystack: str, needle: str) -> int:
        length_needle, length_haystack = len(needle), len(haystack)
        last_needle_index = length_needle - 1

        if length_needle > length_haystack:
            return -1
        
        if haystack == needle:
            return 0

        if length_needle == length_haystack:
            return -1
        
        if needle == haystack:
            return 0

        j = length_needle
        for i, _ in enumerate(haystack):
            current_needle = haystack[i:j]
            if current_needle == needle:
                return i
            j += 1
        
        return -1
