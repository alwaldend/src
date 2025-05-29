/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func zigzagLevelOrder(root *TreeNode) [][]int {
    if root == nil {
        return [][]int{}
    }

    queue, answer := []*TreeNode{root}, [][]int{}
    ltr := true

    for length := len(queue); length != 0; length = len(queue) {
        answerRow := []int{}
        for i := 0; i < length; i++ {
            node := queue[i]
            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
            answerRow = append(answerRow, node.Val)
        }
        if !ltr {
            slices.Reverse(answerRow)
        }
        queue = queue[length:]
        ltr = !ltr
        answer = append(answer, answerRow)
    }

    return answer
}