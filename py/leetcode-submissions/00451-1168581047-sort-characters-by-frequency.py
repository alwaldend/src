class Solution:
    def frequencySort(self, s: str) -> str:
        counter = defaultdict(int)
        for char in s:
            counter[char] += 1
        return "".join(sorted(s, key=lambda val: -(ord(val) + counter[val] * 100)))