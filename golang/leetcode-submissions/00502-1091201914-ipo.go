type Project struct {
    profit, capital int
}

type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func findMaximizedCapital(k int, w int, profits []int, capital []int) int {
    n := len(profits)
    projects := make([]Project, n)
    for i := range profits {
        projects[i] = Project{profits[i], capital[i]}
    }
    sort.Slice(projects, func (i, j int) bool {
        return projects[i].capital < projects[j].capital
    })
    
    q := &IntHeap{}
    heap.Init(q)

    ptr := 0
    for i := 0; i < k; i++ {
        for ptr < n && projects[ptr].capital <= w {
            heap.Push(q, projects[ptr].profit)
            ptr++
        }
        if q.Len() == 0 {
            break
        }
        w += heap.Pop(q).(int)
    }
    return w
}