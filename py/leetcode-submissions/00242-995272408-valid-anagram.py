class Solution:
    def isAnagram(self, s: str, t: str) -> bool:
        s_length, t_length = len(s), len(t)

        if s_length != t_length:
            return False
        
        count = defaultdict(int)
        for i, s_symbol in enumerate(s):
            t_symbol = t[i]
            count[s_symbol] += 1
            count[t_symbol] -= 1

        return not any((i != 0 for i in count.values())) 