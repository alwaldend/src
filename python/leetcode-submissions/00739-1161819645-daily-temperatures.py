class Solution:
    def dailyTemperatures(self, temperatures: List[int]) -> List[int]:
        temps_left = defaultdict(list)
        to_pop = []
        length = len(temperatures)
        ans = [0] * length
        for i in range(length):
            temp = temperatures[i]
            for temp_left, ids in temps_left.items():
                if temp <= temp_left:
                    continue
                for id in ids:
                    ans[id] = i - id
                to_pop.append(temp_left)
            temps_left[temp].append(i)
            for temp in to_pop:
                temps_left.pop(temp)
            to_pop.clear()
        
        return ans