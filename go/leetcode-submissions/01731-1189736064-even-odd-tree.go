/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func isEvenOddTree(root *TreeNode) bool {
    q := []*TreeNode{root}
    for i:= 0; len(q) > 0; i++ {
        var prev *TreeNode
        for _, node := range q {
            q = q[1:]

            if i % 2 == 0 && node.Val % 2 != 1 {
                return false
            }

            if i % 2 == 1 && node.Val % 2 != 0 {
                return false
            }

            if prev != nil && i % 2 == 0 && prev.Val >= node.Val {
                return false
            }

            if prev != nil && i % 2 == 1 && prev.Val <= node.Val {
                return false
            }
           
            prev = node
            if node.Left != nil {
                q = append(q, node.Left)
            }

            if node.Right != nil {
                q = append(q, node.Right)
            }
        }
    }
    return true
}