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
    leftToRight := true

    for length := len(queue); length != 0; length = len(queue) {
        answerRow := []int{}

        for i := 0; i < length; i++ {
            node := queue[0]
            queue = queue[1:]

            answerRow = append(answerRow, node.Val)
            if left := node.Left; left != nil {
                queue = append(queue, left)
            }
            if right := node.Right; right != nil {
                queue = append(queue, right)
            }
        }
        if !leftToRight {
            slices.Reverse(answerRow)
        }
        leftToRight = !leftToRight
        answer = append(answer, answerRow)
    }

    return answer
}