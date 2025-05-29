

func findMinHeightTrees(n int, edges [][]int) []int {
	counts := make([]int, n)
	links := make([]int, n)
	for _, edge := range edges {
		links[edge[0]] ^= edge[1]
		counts[edge[0]]++
		links[edge[1]] ^= edge[0]
		counts[edge[1]]++
	}
	Qu := make([]int, 0)
	dists := make([]int, n)
	for i := 0; i < n; i++ {
		if counts[i] == 1 {
			Qu = append(Qu, i)
		}
	}
	stp := 1
	for len(Qu) > 0 {
		size := len(Qu)
		for j := 0; j < size; j++ {
			tmp := Qu[0]
			Qu = Qu[1:]
			links[links[tmp]] ^= tmp
			counts[links[tmp]]--
			if counts[links[tmp]] == 1 {
				dists[links[tmp]] = int(math.Max(float64(stp), float64(dists[links[tmp]])))
				Qu = append(Qu, links[tmp])
			}
		}
		stp++
	}
	maxDist := 0
	for _, dist := range dists {
		if dist > maxDist {
			maxDist = dist
		}
	}
	res := make([]int, 0)
	for i, dist := range dists {
		if dist == maxDist {
			res = append(res, i)
		}
	}
	return res
}