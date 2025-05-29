/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func findMode(root *TreeNode) []int {
    if root == nil {
        return []int{}
    }
    counter := map[int]int{}
    queue := []*TreeNode{root}
    for length := len(queue); length > 0; length = len(queue) {
        for i := 0; i < length; i++ {
            node := queue[i]
            if node == nil {
                continue
            }
            if val, ok := counter[node.Val]; ok {
                counter[node.Val] = val + 1
            } else {
                counter[node.Val] = 1
            }
            queue = append(queue, node.Left, node.Right)
        }
        queue = queue[length:]
    }
    maxCount, answer := 0, []int{}
    for key, count := range counter {
        if count < maxCount {
            continue
        }
        if count > maxCount {
            answer = answer[:0]
            maxCount = count
        }
        answer = append(answer, key)
    }
    return answer
}