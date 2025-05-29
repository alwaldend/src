/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func maxPathSum(root *TreeNode) int {
    ans := -1 << 63
    maxPath(root, &ans)
    return ans
}

func maxPath(root *TreeNode, ans *int) int {
    if root == nil {
        return 0
    }
    leftPathSum := maxPath(root.Left, ans)
    rightPathSum := maxPath(root.Right, ans)
    *ans = max(*ans, leftPathSum + rightPathSum + root.Val)
    return max(max(leftPathSum+root.Val, rightPathSum+root.Val), 0)
}

func max(a int, b int) int {
    if (a >= b) {
        return a
    }
    return b
}
