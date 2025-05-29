/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func averageOfSubtree(root *TreeNode) int {
    _, _, answer := getAverage(root)
    return answer
}

func getAverage(root *TreeNode) (int, int, int) {
    if root == nil {
        return 0, 0, 0
    }
    sum, count, validNodes := root.Val, 1, 0
    leftSum, leftCount, leftValidNodes := getAverage(root.Left)
    rightSum, rightCount, rightValidNodes := getAverage(root.Right)
    validNodes += leftValidNodes + rightValidNodes
    sum += leftSum + rightSum
    count += leftCount + rightCount
    if sum / count == root.Val {
        validNodes += 1
    }
    return sum, count, validNodes
}