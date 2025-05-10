/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

type queueItem struct {
    node *TreeNode
    sum int
}
func sumNumbers(root *TreeNode) int {
    if root == nil {
        return 0
    }
    queue := []queueItem{queueItem{root, 0}}
    sum := 0
    for length := len(queue); length > 0; length = len(queue) {
        for i := 0; i < length; i++ {
            item := queue[i]
            if item.node == nil {
                continue
            }
            left, right, newVal := item.node.Left, item.node.Right, item.node.Val + item.sum * 10
            if left == nil && right == nil {
                sum += newVal
            } else {
                queue = append(queue, queueItem{left, newVal}, queueItem{right, newVal})
            }
        }
        queue = queue[length:]
    }
    return sum
}