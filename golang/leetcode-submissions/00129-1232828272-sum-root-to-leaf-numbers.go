/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func sumNumbers(root *TreeNode) int {
    return dfs(root, 0)
}

func dfs(node *TreeNode, num int) int {
    if node == nil {
        return 0
    }
    if node.Left == nil && node.Right == nil {
        return num * 10 + node.Val
    }
    return dfs(node.Left, num * 10 + node.Val) + dfs(node.Right, num * 10 + node.Val)
}