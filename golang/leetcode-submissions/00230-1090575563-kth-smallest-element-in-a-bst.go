/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func kthSmallest(root *TreeNode, k int) int {
    heap := binaryheap.NewWithIntComparator()
    queue := []*TreeNode{root}

    for length := len(queue); length > 0; length = len(queue) {
        for i := 0; i < length; i++ {
            node := queue[i]
            if node == nil {
                continue
            }
            heap.Push(node.Val)
            queue = append(queue, node.Left, node.Right)
        }
        queue = queue[length:]
    }
    var answer any
    for i := 0; i < k; i++ {
        answer, _ = heap.Pop()
    }
    return answer.(int)
}
