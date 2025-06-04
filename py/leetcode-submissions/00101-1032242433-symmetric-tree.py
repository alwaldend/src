# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def isSymmetric(self, root: Optional[TreeNode]) -> bool:
        comparison_stack = [(root.left, root.right)]
        while comparison_stack:
            left, right = comparison_stack.pop()
            
            if not left and not right:
                continue
            if not left or not right or left.val != right.val:
                return False

            comparison_stack.extend(((left.left, right.right), (left.right, right.left)))
        
        return True

