class Solution:
    def simplifyPath(self, path: str) -> str:
        canonical = []
        for directory in path.split("/"):
            if not directory or directory == ".":
                continue
            
            if directory != "..":
                canonical.append(directory)
                continue
            
            if len(canonical):
                canonical.pop()

        return "/" + "/".join(canonical)