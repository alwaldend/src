# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class FindElements:

    def __init__(self, root: Optional[TreeNode]):
        root.val = 0
        nodes = [root]
        self.values = set((0, )) 
        while nodes:
            node = nodes.pop()
            if node.left:
                val = 2 * node.val + 1
                node.left.val = val
                nodes.append(node.left)
                self.values.add(val)
            if node.right:
                val = 2 * node.val + 2
                node.right.val = val
                nodes.append(node.right)
                self.values.add(val)

    def find(self, target: int) -> bool:
        return target in self.values

# Your FindElements object will be instantiated and called as such:
# obj = FindElements(root)
# param_1 = obj.find(target)