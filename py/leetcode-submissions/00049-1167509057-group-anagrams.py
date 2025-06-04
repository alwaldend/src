class Solution:
    def groupAnagrams(self, strs: List[str]) -> List[List[str]]:
        anagrams = defaultdict(list)
        for anagram in strs:
            anagrams[tuple(sorted(anagram))].append(anagram)
        return anagrams.values()