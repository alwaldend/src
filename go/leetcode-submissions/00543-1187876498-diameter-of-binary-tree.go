var maxD int

func diameterOfBinaryTree(root *TreeNode) int {
    maxD = 0
    find(root)
    return maxD
}

func find(root *TreeNode) int {
    if root == nil {
        return 0
    }
    left := find(root.Left)
    right := find(root.Right)
    localMax := left + right
    maxD = max(maxD, localMax)
    return max(left, right) + 1
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}