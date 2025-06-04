class Solution:
    
    def longestCommonPrefix(self, strs: List[str]) -> str:
        if len(strs) == 1:
            return strs[0]

        min_length = min([len(string) for string in strs])
        for i in range(min_length, -1, -1):
            current = strs[0][0:i+1]
            for string in strs[1:]:
                if string[0:i+1] != current:
                    break
            else:
                return current

        return ""    