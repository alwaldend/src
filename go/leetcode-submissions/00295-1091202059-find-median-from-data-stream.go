type MaxHeap []int
func (m MaxHeap) Len() int { return len(m) }
func (m MaxHeap) Less(i, j int) bool { return m[i] > m[j] }
func (m MaxHeap) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m *MaxHeap) Pop() interface{} {
    v := (*m)[len(*m)-1]
    *m = (*m)[:len(*m)-1]
    return v
}
func (m *MaxHeap) Push(v interface{}) { *m = append(*m, v.(int)) }
func (m MaxHeap) Top() int { return m[0] }

type MinHeap []int
func (m MinHeap) Len() int { return len(m) }
func (m MinHeap) Less(i, j int) bool { return m[i] < m[j] }
func (m MinHeap) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m *MinHeap) Pop() interface{} {
    v := (*m)[len(*m)-1]
    *m = (*m)[:len(*m)-1]
    return v
}
func (m *MinHeap) Push(v interface{}) { *m = append(*m, v.(int)) }
func (m MinHeap) Top() int { return m[0] }

type MedianFinder struct {
    left MaxHeap
    right MinHeap    
}

func Constructor() MedianFinder {
    return MedianFinder{}    
}

func (mf *MedianFinder) AddNum(num int)  {
    if len(mf.left) + len(mf.right) == 0 {
        heap.Push(&(mf.left), num)
        return
    }
    for {
        if len(mf.left) < len(mf.right) {
            if num <= mf.right.Top() {
                heap.Push(&(mf.left), num)
                return
            } else {
                v := heap.Pop(&(mf.right))
                heap.Push(&(mf.left), v)
            }
        } else {
            if num >= mf.left.Top() {
                heap.Push(&(mf.right), num)
                return
            } else {
                v := heap.Pop(&(mf.left))
                heap.Push(&(mf.right), v)
            }
        }
    }
}

func (mf *MedianFinder) FindMedian() float64 {
    if len(mf.left) == len(mf.right) {
        return float64(mf.left.Top() + mf.right.Top()) / 2.0
    } else if len(mf.left) > len(mf.right) {
        return float64(mf.left.Top())
    } else {
        return float64(mf.right.Top())
    }   
}