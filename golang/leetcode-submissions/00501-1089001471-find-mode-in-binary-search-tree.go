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
    dfs(root, counter)
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

func dfs(root *TreeNode, counter map[int]int) {
    if root == nil {
        return
    }
    if _, ok := counter[root.Val]; ok {
        counter[root.Val] += 1
    } else {
        counter[root.Val] = 1
    }
    dfs(root.Left, counter)
    dfs(root.Right, counter)
}