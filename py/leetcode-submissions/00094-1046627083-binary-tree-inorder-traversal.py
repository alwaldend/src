# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def inorderTraversal(self, root: Optional[TreeNode]) -> List[int]:
        def get_nodes(root: TreeNode) -> Generator[None, None, TreeNode]:
            if not root:
                return

            yield from get_nodes(root.left)
            yield root.val
            yield from get_nodes(root.right)
        
        return tuple(get_nodes(root))