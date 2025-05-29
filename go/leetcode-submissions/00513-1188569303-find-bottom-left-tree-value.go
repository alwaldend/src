/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func findBottomLeftValue(root *TreeNode) int {
    _, val := getBottomLeft(root)
    return val
}

func getBottomLeft(root *TreeNode) (int, int) {
    left, right := root.Left, root.Right
    if right == nil && left == nil {
        return 0, root.Val
    }
    if right == nil {
        depth, val := getBottomLeft(root.Left)
        return 1 + depth, val
    }
    if left == nil {
        depth, val := getBottomLeft(root.Right)
        return 1 + depth, val
    }
    leftDepth, leftVal := getBottomLeft(root.Left)
    rightDepth, rightVal := getBottomLeft(root.Right)
    if rightDepth > leftDepth {
        return 1 + rightDepth, rightVal
    }
    return 1 + leftDepth, leftVal
}