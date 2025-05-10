from collections import deque

class Solution:
  def connect(self, root: 'Node') -> 'Node':
    # Edge case - If the root is None, then return None
    if root is None:
        return None
    
    # Create a queue and enqueue the root node
    q = deque([root])
    
    # Traverse the tree level by level
    while q:
        
        # Get the number of nodes of the current level
        level_size = len(q)
        
        # Process the nodes of the current level
        for i in range(level_size):
            
            # Dequeue a node from the front of the queue
            node = q.popleft()
            
            # Assign the next pointer of the node
            if i < level_size - 1:
                node.next = q[0]
            
            # Enqueue the children of the node (if any)
            if node.left is not None:
                q.append(node.left)
            if node.right is not None:
                q.append(node.right)
    
    # Return the root node
    return root
