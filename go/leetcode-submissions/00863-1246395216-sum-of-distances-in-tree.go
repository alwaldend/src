
func sumOfDistancesInTree(n int, edges [][]int) []int {
	graph := make(map[int][]int)
	count := make([]int, n)
	res := make([]int, n)

	for _, edge := range edges {
		u, v := edge[0], edge[1]
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	var dfs1 func(cur, parent int)
	dfs1 = func(cur, parent int) {
		count[cur] = 1
		for _, child := range graph[cur] {
			if child != parent {
				dfs1(child, cur)
				count[cur] += count[child]
				res[cur] += res[child] + count[child]
			}
		}
	}

	var dfs2 func(cur, parent int)
	dfs2 = func(cur, parent int) {
		for _, child := range graph[cur] {
			if child != parent {
				res[child] = res[cur] + n - 2*count[child]
				dfs2(child, cur)
			}
		}
	}

	dfs1(0, -1)
	dfs2(0, -1)

	return res
}
