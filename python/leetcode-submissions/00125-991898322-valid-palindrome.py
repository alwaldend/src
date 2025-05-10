class Solution:
    def isPalindrome(self, s: str) -> bool:
        length = len(s)

        # if the sring has only one symbol and it is alphanumeric, it is a palyndrom
        if s.isalnum() and length == 1:
            return True 
        
        i, j = 0, length - 1
        # iterate from the start and from the end using two indexes:
        # - if one the symbols is not alphanumeric, move the corresponding index
        # - if symbols are alphanumeric and not equal, return False
        # - move indexes
        while i < j:
            symbol_start, symbol_end = s[i].lower(), s[j].lower()
            if not symbol_start.isalnum():
                i += 1
                continue
            if not symbol_end.isalnum():
                j -= 1
                continue
            
            if symbol_start != symbol_end:
                return False
            
            i += 1
            j -= 1
        
        return True