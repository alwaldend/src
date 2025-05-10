# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def isBalanced(self, root: Optional[TreeNode]) -> bool:
        def get_height(node: TreeNode) -> int:
            if not node:
                return 0
            
            h_left = get_height(node.left)
            if h_left < 0:
                return -1
            h_right = get_height(node.right)
            if h_right < 0 or abs(h_left - h_right) > 1:
                return -1

            return max(h_left, h_right) + 1
        
        return get_height(root) >= 0