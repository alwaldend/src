/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
type queueItem struct  {
    node *TreeNode
    max *int
    min *int
}

func isValidBST(root *TreeNode) bool {
    queue := []queueItem{queueItem{root, nil, nil}}
    for length := len(queue); length > 0; length = len(queue) {
        for i := 0; i < length; i++ {
            item := queue[i]
            if item.node == nil {
                continue
            }
            node := item.node
            if (item.max != nil && node.Val >= *item.max) || (item.min != nil && node.Val <= *item.min) {
                return false
            }
            queue = append(
                queue, 
                queueItem{node.Left, &node.Val, item.min}, 
                queueItem{node.Right, item.max, &node.Val},
            )
        }
        queue = queue[length:]
    }
    return true
}