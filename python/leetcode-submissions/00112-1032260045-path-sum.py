# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def hasPathSum(self, root: Optional[TreeNode], targetSum: int) -> bool:
        node_stack = [(root, 0)]
        while node_stack:
            node, curr_sum = node_stack.pop()
            if not node:
                continue
            
            new_sum = curr_sum + node.val
            if new_sum == targetSum and not node.left and not node.right:
                return True
            
            node_stack.extend(((node.left, new_sum), (node.right, new_sum)))
                
        return False