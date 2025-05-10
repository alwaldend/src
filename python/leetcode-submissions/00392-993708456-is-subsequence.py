class Solution:
    # create two pointers, original and sub, both are zero
    # while either of those pointers have not reached the end:
    # - if original symbol is equal to the sub, move both pointer to the right
    # - if not, move original pointer to the right
    # if sub pointer reached the end, return True, otherwise False
    def isSubsequence(self, s: str, t: str) -> bool:
        original, sub = 0, 0
        original_length, sub_length = len(t), len(s)

        while original < original_length and sub < sub_length:
            if s[sub] == t[original]:
                original += 1
                sub += 1
                continue
            
            original += 1

        return sub == sub_length