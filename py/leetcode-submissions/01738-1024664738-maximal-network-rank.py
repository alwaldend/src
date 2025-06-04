class Solution:
    def maximalNetworkRank(self, n: int, roads: List[List[int]]) -> int:
        city_roads = [set() for _ in range(n)]

        for city_one, city_two in roads:
            city_roads[city_one].add(city_two)
            city_roads[city_two].add(city_one)
        
        max_rank = 0

        for city_one in range(n):
            for city_two in range(city_one + 1, n):
                rank = len(city_roads[city_one]) + len(city_roads[city_two])
                if city_one in city_roads[city_two]:
                    rank -= 1
                
                if rank > max_rank:
                    max_rank = rank
        
        return max_rank