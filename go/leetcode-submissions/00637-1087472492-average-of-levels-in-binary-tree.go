/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func averageOfLevels(root *TreeNode) []float64 {
    if root == nil {
        return []float64{}
    }
    queue, answer := []*TreeNode{root}, []float64{}

    for length := len(queue); length != 0; length = len(queue) {
        av := 0
        for i := 0; i < length; i++ {
            node := queue[0]
            queue = queue[1:]
            av += node.Val
            if left := node.Left; left != nil {
                queue = append(queue, left)
            }
            if right := node.Right; right != nil {
                queue = append(queue, right)
            }
        }
        answer = append(answer, float64(av) / float64(length))
    }

    return answer
}