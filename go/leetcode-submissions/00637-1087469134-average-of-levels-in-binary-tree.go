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
    qCur, qNext, answer := []*TreeNode{root}, []*TreeNode{}, []float64{}

    for len(qCur) != 0 {
        av := 0
        for _, node := range qCur {
            av += node.Val
            if left := node.Left; left != nil {
                qNext = append(qNext, left)
            }
            if right := node.Right; right != nil {
                qNext = append(qNext, right)
            }
        }
        answer = append(answer, float64(av) / float64(len(qCur)))
        qCur = qCur[:0]
        qCur, qNext = qNext, qCur
    }

    return answer
}