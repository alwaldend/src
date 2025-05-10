# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def countNodes(self, root: Optional[TreeNode]) -> int:
        def left_height(root: TreeNode) -> int:
            return 0 if not root else 1 + left_height(root.left)
        
        def right_height(root: TreeNode) -> int:
            return 0 if not root else 1 + right_height(root.right)
        
        left, right = left_height(root), right_height(root)
        if left > right:
            return 1 + self.countNodes(root.left) + self.countNodes(root.right)
        
        return 2**left - 1

        