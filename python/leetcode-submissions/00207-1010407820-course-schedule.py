class Solution:
    def canFinish(self, numCourses: int, prerequisites: List[List[int]]) -> bool:
        courses = defaultdict(set)

        for course, prereq in prerequisites:
            courses[course].add(prereq)

        if not courses:
            return True
        
        stack = set()

        @cache
        def check(course: int) -> bool:
            if course not in courses:
                return True
            
            if course in stack:
                return False

            prereqs = courses[course]

            if prereqs & stack:
                return False

            stack.add(course)

            for prereq in prereqs:
                if not check(prereq):
                    return False
    
            stack.remove(course)
            
            return True

        for course in courses:
            if not check(course):
                return False

        return True