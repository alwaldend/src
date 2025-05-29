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

    queue, answer := []*TreeNode{root}, []int{}

    for length := len(queue); length != 0; length = len(queue) {
        for i := 0; i < length; i++ {
            node := queue[0]
            queue = queue[1:]
            if i == length - 1 {
                answer = append(answer, node.Val)
            }
            if left := node.Left; left != nil {
                queue = append(queue, left)
            }
            if right := node.Right; right != nil {
                queue = append(queue, right)
            }
        }
    }

    return answer
} 