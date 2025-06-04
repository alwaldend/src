class Solution:
    def calcEquation(self, equations: List[List[str]], values: List[float], queries: List[List[str]]) -> List[float]:
        equations_count = len(equations)
        queries_count = len(queries) 
        edges = defaultdict(dict)
        seen_nodes = set()
        invalid = -1.0

        def dfs(val1: str, val2: str) -> float:
            if val1 == val2:
                return 1.0

            if val2 in edges[val1]:
                return edges[val1][val2]

            if val1 in seen_nodes:
                return invalid

            seen_nodes.add(val1)

            for connected_node in edges[val1]:
                new_edge_res = dfs(connected_node, val2)
                if new_edge_res == invalid:
                    continue
                new_edge_res *= edges[val1][connected_node]
                edges[val1][val2] = new_edge_res
                return new_edge_res
            
            return invalid

        for i in range(equations_count):
            (val1, val2), res = equations[i], values[i]
            edges[val1][val2], edges[val2][val1] = res, 1 / res
        
        for i in range(queries_count):
            res, (val1, val2) = -1, queries[i]
            if val1 in edges and val2 in edges:
                seen_nodes.clear()
                res = dfs(val1, val2)
            queries[i] = res

        return queries