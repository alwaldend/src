"""
# Definition for a Node.
class Node:
    def __init__(self, val = 0, neighbors = None):
        self.val = val
        self.neighbors = neighbors if neighbors is not None else []
"""

class Solution:
    def cloneGraph(self, node: 'Node') -> 'Node':
        old_to_new = {}
        
        def clone(node: Node) -> Node:
            if node in old_to_new:
                return old_to_new[node]
            
            new_node = Node(node.val)
            old_to_new[node] = new_node
            new_node.neighbors = [clone(neighbor) for neighbor in node.neighbors]

            return new_node

        return clone(node) if node else node