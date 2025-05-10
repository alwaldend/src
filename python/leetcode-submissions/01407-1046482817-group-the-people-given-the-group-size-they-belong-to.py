class Solution:
    def groupThePeople(self, groupSizes: List[int]) -> List[List[int]]:
        answer = []
        cur_groups = defaultdict(list)
        
        for i in range(len(groupSizes)):
            size = groupSizes[i]
            group = cur_groups[size]
            group.append(i)
            if len(group) == size:
                answer.append(group)
                del cur_groups[size]
        
        return answer