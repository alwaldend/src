"""
# Definition for a Node.
class Node:
    def __init__(self, x: int, next: 'Node' = None, random: 'Node' = None):
        self.val = int(x)
        self.next = next
        self.random = random
"""

class Solution:
    def copyRandomList(self, head: 'Optional[Node]') -> 'Optional[Node]':
        
        copied = {None: None}

        def copy_node(node: Node) -> Node:
            if node in copied:
                return copied[node]
            
            new_node = Node(node.val)
            copied[node] = new_node
            new_node.next = copy_node(node.next)
            new_node.random = copy_node(node.random)
            return new_node
        
        return copy_node(head)
