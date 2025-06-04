# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def getMinimumDifference(self, root: Optional[TreeNode]) -> int:
        values = []
        def fall(node: TreeNode) -> None:
            if not node:
                return
            values.append(node.val)
            if node.left:
                fall(node.left)
            if node.right:
                fall(node.right)
        
        fall(root)
        values.sort()
        min_diff = abs(values[1] - values[0])
        for i in range(1, len(values) - 1):
            min_diff = min(min_diff, abs(values[i] - values[i+1]))
        
        return min_diff