/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    isAncestorMap := map[int][2]bool{}
    dfs(root, p, q, isAncestorMap)
    isAncestor := func (node *TreeNode) bool {
        if node == nil {
            return false
        }
        hasNodes := isAncestorMap[node.Val]
        return hasNodes[0] && hasNodes[1]
    }

    for {
        left, right := root.Left, root.Right
        if isAncestor(left) {
            root = left
        } else if isAncestor(right) {
            root = right
        } else {
            break
        }
    }
    return root
}

func dfs(root, p, q *TreeNode, isAncestorMap map[int][2]bool) (bool, bool)  {
    if root == nil {
        return false, false
    }
    hasPLeft, hasQLeft := dfs(root.Left, p, q, isAncestorMap)
    hasPRight, hasQRight := dfs(root.Right, p, q, isAncestorMap)
    hasP, hasQ := hasPLeft || hasPRight, hasQLeft || hasQRight
    if root == p {
        hasP = true
    }
    if root == q {
        hasQ = true
    }
    isAncestorMap[root.Val] = [2]bool{hasP, hasQ}
    return hasP, hasQ
}