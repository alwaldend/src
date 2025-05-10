# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def largestValues(self, root: Optional[TreeNode]) -> List[int]:
        if not root:
            return []

        largest = []
        cur, nxt = [root], []
        cur_val = None
        while cur or nxt:
            while cur:
                node = cur.pop()
                if cur_val is None or node.val > cur_val:
                    cur_val = node.val
                if node.left:
                    nxt.append(node.left)
                if node.right:
                    nxt.append(node.right)
            largest.append(cur_val)
            cur.clear()
            cur_val = None
            cur, nxt = nxt, cur
        return largest