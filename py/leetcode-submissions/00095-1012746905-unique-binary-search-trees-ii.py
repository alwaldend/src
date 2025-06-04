class Solution:
    def generateTrees(self, n: int) -> List[Optional[TreeNode]]:
        @cache
        def generate_trees(l, r):
            return [None] if l > r else [
                TreeNode(val, left, right)
                for val in range(l, r + 1)
                for left in generate_trees(l, val - 1)
                for right in generate_trees(val + 1, r)
            ]
        
        return generate_trees(1, n)