/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func rightSideView(root *TreeNode) []int {
    if root == nil {
        return []int{}
    }

    q_cur, q_next, answer := []*TreeNode{root}, []*TreeNode{}, []int{}

    for length := len(q_cur); length > 0; length = len(q_cur) {
        for _, node := range q_cur {
            if left := node.Left; left != nil {
                q_next = append(q_next, left)
            }
            if right := node.Right; right != nil {
                q_next = append(q_next, right)
            }
        }

        answer = append(answer, q_cur[length - 1].Val)
        q_cur = q_cur[:0]
        q_cur, q_next = q_next, q_cur
    }

    return answer
} 