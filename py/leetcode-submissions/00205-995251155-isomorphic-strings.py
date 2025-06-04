class Solution:
    def isIsomorphic(self, s: str, t: str) -> bool:
        s_length, t_length = len(s), len(t)
        if s_length != t_length:
            return False
        
        s_to_t = {}
        t_to_s = {}
        for i, s_symbol in enumerate(s):
            t_symbol = t[i]

            if s_symbol not in s_to_t and t_symbol not in t_to_s:
                s_to_t[s_symbol] = t_symbol
                t_to_s[t_symbol] = s_symbol
            elif s_symbol in s_to_t and s_to_t[s_symbol] == t_symbol:
                continue
            else:
                return False
                
        return True