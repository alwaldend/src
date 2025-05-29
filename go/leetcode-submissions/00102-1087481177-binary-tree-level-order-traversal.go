/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return [][]int{}
    }

    stack, answer := []*TreeNode{root}, [][]int{}

    for length := len(stack); length != 0; length = len(stack) {
        answerRow := []int{}
        for i := 0; i < length; i++ {
            node := stack[0]
            answerRow = append(answerRow, node.Val)
            stack = stack[1:]
            if left := node.Left; left != nil {
                stack = append(stack, node.Left)
            }
            if right := node.Right; right != nil {
                stack = append(stack, node.Right)
            }
        }
        answer = append(answer, answerRow)
    }

    return answer
}