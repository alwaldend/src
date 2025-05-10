# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def sumOfLeftLeaves(self, root: Optional[TreeNode]) -> int:
        nodes = set([root])
        result = 0

        while nodes:
            node = nodes.pop()
            
            if node.left and not node.left.left and not node.left.right:
                result += node.left.val
            elif node.left:
                nodes.add(node.left)
            
            if node.right:
                nodes.add(node.right)
            
        return result