# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def rightSideView(self, root: Optional[TreeNode]) -> List[int]:
        if not root:
            return []
        cur, nxt = [root], []
        answer = []
        
        while cur:
            last_node = cur.pop()
            if last_node.right:
                nxt.append(last_node.right)
            if last_node.left:
                nxt.append(last_node.left)
            answer.append(last_node.val)
            
            while cur:
                node = cur.pop()
                if node.right:
                    nxt.append(node.right)
                if node.left:
                    nxt.append(node.left)
            
            cur.clear()
            nxt.reverse()
            cur, nxt = nxt, cur

        return answer