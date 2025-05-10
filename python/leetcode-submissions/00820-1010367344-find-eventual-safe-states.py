class Solution:
    def eventualSafeNodes(self, graph: List[List[int]]) -> List[int]:
        length = len(graph)
        indegree = [0] * length
        adj = [[] for _ in range(length)]

        for i in range(length):
            for edge in graph[i]:
                adj[edge].append(i)
                indegree[i] += 1

        q = deque()
        # Push all the nodes with indegree zero in the queue.
        for i in range(length):
            if indegree[i] == 0:
                q.append(i)

        safe = [False] * length
        while q:
            node = q.popleft()
            safe[node] = True

            for neighbor in adj[node]:
                # Delete the edge "node -> neighbor".
                indegree[neighbor] -= 1
                if indegree[neighbor] == 0:
                    q.append(neighbor)

        safe_nodes = []
        for i in range(length):
            if safe[i]:
                safe_nodes.append(i)

        return safe_nodes