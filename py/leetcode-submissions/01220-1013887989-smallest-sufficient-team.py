class Solution:
    def smallestSufficientTeam(self, req_skills: List[str], people: List[List[str]]) -> List[int]:
        req_skills_count = len(req_skills)
        skill_to_people = defaultdict(set)
        for i, person in enumerate(people):
            for skill in person:
                skill_to_people[skill].add(i)

        current = set()
        
        def backtrack(skill: int) -> Tuple[int]:
            if skill == req_skills_count:
                return tuple(current)
            
            suff_team = None
            people_with_skill = skill_to_people[req_skills[skill]] 
            
            if current & people_with_skill:
                return backtrack(skill + 1)
        
            for person in people_with_skill:
                current.add(person)
                new_team = backtrack(skill + 1)
                current.remove(person)

                if suff_team is None or len(new_team) < len(suff_team):
                    suff_team = new_team

            return suff_team

        return backtrack(0)
        
