# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def maxDepth(self, root: Optional[TreeNode], count: int = 0) -> int:
        if not root:
            return 0
        
        if not root.left and not root.right:
            return count + 1

        return max(self.maxDepth(root.left, count + 1), self.maxDepth(root.right, count + 1))
        