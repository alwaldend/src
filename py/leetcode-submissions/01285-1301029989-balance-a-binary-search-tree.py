# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def balanceBST(self, root: TreeNode) -> TreeNode:
        self.sortedArr = []
        self.inorderTraverse(root)
        return self.sortedArrayToBST(0, len(self.sortedArr) - 1)

    def inorderTraverse(self, root: TreeNode) -> None:
        if not root:
            return
        self.inorderTraverse(root.left)
        self.sortedArr.append(root)
        self.inorderTraverse(root.right)

    def sortedArrayToBST(self, start: int, end: int) -> TreeNode:
        if start > end:
            return None
        mid = (start + end) // 2
        root = self.sortedArr[mid]
        root.left = self.sortedArrayToBST(start, mid - 1)
        root.right = self.sortedArrayToBST(mid + 1, end)
        return root        