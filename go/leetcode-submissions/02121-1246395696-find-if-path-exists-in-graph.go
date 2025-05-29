type UnionFind struct {
    parent []int
    rank   []int
}

func NewUnionFind(n int) *UnionFind {
    parent := make([]int, n)
    rank := make([]int, n)
    for i := 0; i < n; i++ {
        parent[i] = i
    }
    return &UnionFind{parent, rank}
}

func (uf *UnionFind) Find(u int) int {
    if uf.parent[u] != u {
        uf.parent[u] = uf.Find(uf.parent[u])
    }
    return uf.parent[u]
}

func (uf *UnionFind) UnionByRank(u, v int) {
    i := uf.Find(u)
    j := uf.Find(v)
    if i == j {
        return
    }
    if uf.rank[i] < uf.rank[j] {
        uf.parent[i] = j
    } else if uf.rank[i] > uf.rank[j] {
        uf.parent[j] = i
    } else {
        uf.parent[i] = j
        uf.rank[j]++
    }
}

func validPath(n int, edges [][]int, source int, destination int) bool {
    uf := NewUnionFind(n)
    for _, edge := range edges {
        u, v := edge[0], edge[1]
        uf.UnionByRank(u, v)
    }
    return uf.Find(source) == uf.Find(destination)
}
