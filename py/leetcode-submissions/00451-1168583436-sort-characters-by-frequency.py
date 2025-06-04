class Solution:
    def frequencySort(self, s: str) -> str:
        counter = {}
        for char in s:
            if char not in counter:
                counter[char] = ord(char)
            counter[char] += 100
        return "".join(sorted(s, key=lambda val: -counter[val]))