class Solution:
    def destCity(self, paths: List[List[str]]) -> str:
        paths_from, paths_to = set(), set()
        for path_from, path_to in paths:
            paths_from.add(path_from)
            paths_to.add(path_to)
        
        return (paths_to - paths_from).pop()