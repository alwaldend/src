class Solution:
    def canFinish(self, numCourses: int, prerequisites: List[List[int]]) -> bool:
        
        nodes = [[] for _ in range(numCourses)]
        for target, required in prerequisites:
            nodes[target].append(required)

        visited = set()

        @cache
        def dfs(course: int) -> bool:
            if course in visited:
                return False
    
            edges = nodes[course]
            if not edges:
                return True
            
            visited.add(course)
            return all(dfs(edge) for edge in edges)

        for course in range(numCourses):
            if not dfs(course):
                return False
            visited.clear()
        
        return True