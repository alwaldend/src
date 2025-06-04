# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def pseudoPalindromicPaths (self, root: Optional[TreeNode]) -> int:
        if not root:
            return 0
        
        num_count = defaultdict(int)
        perm_count = 0

        def traverse(node: TreeNode) -> int:
            if not node:
                return 0
            val, left, right = node.val, node.left, node.right
            num_count[val] += 1
            res = 0
            if not left and not right:
                non_even = 0
                for num in num_count.values():
                    if num % 2 == 0:
                        continue
                    non_even += 1
                    if non_even > 1:
                        break
                else:
                    res += 1
            res += traverse(left) + traverse(right)
            num_count[val] = max(0, num_count[val] - 1)
            return res 

        return traverse(root)