class Node:
    def __init__(self, val: int = 0, prev_node: 'Node' = None, next_node: 'Node' = None):
        self.val, self.prev, self.next = val, prev_node, next_node

class MyLinkedList:

    def __init__(self):
        self._head, self._tail = Node(float("-inf")), Node(float("inf"))
        self._head.next, self._tail.prev = self._tail, self._head
        self._length = 0

    def _getNode(self, index: int) -> Node:
        if not 0 <= index < self._length:
            return None
        if index == 0:
            return self._head.next
        if index == self._length - 1:
            return self._tail.prev

        cur_idx, cur_node = 0, self._head.next
        while cur_idx < index:
            cur_node = cur_node.next
            cur_idx += 1

        return cur_node
    
    def get(self, index: int) -> int:
        node = self._getNode(index)
        return node.val if node else -1

    def addAtHead(self, val: int) -> None:
        old_first = self._head.next
        new_node = Node(val, self._head, old_first)
        self._head.next, old_first.prev = new_node, new_node
        self._length += 1

    def addAtTail(self, val: int) -> None:
        old_tail = self._tail.prev
        new_node = Node(val, old_tail, self._tail)
        self._tail.prev, old_tail.next = new_node, new_node
        self._length += 1

    def addAtIndex(self, index: int, val: int) -> None:
        if not 0 <= index <= self._length:
            return
        if index == 0:
            self.addAtHead(val)
            return
        if index == self._length:
            self.addAtTail(val)
            return
        
        target_node = self._getNode(index)
        prev_node = target_node.prev
        new_node = Node(val, prev_node, target_node)
        prev_node.next, target_node.prev = new_node, new_node
        self._length += 1

    def deleteAtIndex(self, index: int) -> None:
        if not 0 <= index < self._length:
            return
        
        target_node = self._getNode(index)
        old_prev, old_next = target_node.prev, target_node.next
        old_prev.next, old_next.prev = old_next, old_prev
        self._length -= 1

# Your MyLinkedList object will be instantiated and called as such:
# obj = MyLinkedList()
# param_1 = obj.get(index)
# obj.addAtHead(val)
# obj.addAtTail(val)
# obj.addAtIndex(index,val)
# obj.deleteAtIndex(index)