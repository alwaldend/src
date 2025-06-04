class Solution:
    def groupAnagrams(self, strs: List[str]) -> List[List[str]]:
        anagrams: Dict[str, List[str]] = defaultdict(list)

        for string in strs:
            anagrams["".join(sorted(string))].append(string)
        
        return anagrams.values()

