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
            left, right, val := node.Left, node.Right, node.Val
            if _, ok := counter[val]; ok {
                counter[val] +=  1
            } else {
                counter[val] = 1
            }
            queue = append(queue, left, right)
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