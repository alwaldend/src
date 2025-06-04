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
        
        q_cur, q_next, answer = deque(), deque(), []
        q_cur.append(root)

        while q_cur:
            last_val = None
            
            while q_cur:
                node = q_cur.popleft()
                last_val = node.val
                if node.left:
                    q_next.append(node.left)
                if node.right:
                    q_next.append(node.right)
            
            q_cur, q_next = q_next, q_cur
            answer.append(last_val)
    
        return answer