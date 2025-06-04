class Solution:
    def isIsomorphic(self, s: str, t: str) -> bool:
        return list(map(s.index, s)) == list(map(t.index, t))