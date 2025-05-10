/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func sumNumbers(root *TreeNode) int {
    return getSum(root, 0)
}

func getSum(root *TreeNode, parentVal int) int {
    if root == nil {
        return 0
    }
    newVal := root.Val + parentVal * 10
    if root.Left == nil && root.Right == nil {
        return newVal
    }
    return getSum(root.Left, newVal) + getSum(root.Right, newVal)
}