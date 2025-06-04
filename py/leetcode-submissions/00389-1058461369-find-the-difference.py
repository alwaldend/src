class Solution:
    def findTheDifference(self, s: str, t: str) -> str:
        return chr(reduce(operator.xor, (ord(char) for char in chain(s, t))))