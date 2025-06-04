# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def getMinimumDifference(self, root: Optional[TreeNode]) -> int:
        
        def min_node_diff(node: TreeNode, low: int, high: int) -> int:
            if not node:
                return high - low

            return min(min_node_diff(node.left, low, node.val),
                       min_node_diff(node.right, node.val, high))
    
        return min_node_diff(root, -maxsize, maxsize)