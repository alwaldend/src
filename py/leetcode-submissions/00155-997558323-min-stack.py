from sortedcontainers import sortedset

class MinStack:
    class Node:
        def __init__(self, val: int, prev: 'Node'):
            self.val = val
            self.prev = prev


    def __init__(self):
        self.stack: List[Node] = []
        self.min_node = None

    def push(self, val: int) -> None:
        new_node = self.Node(val, None)
        self.stack.append(new_node)
        if self.min_node is None or val <= self.min_node.val:
            new_node.prev = self.min_node
            self.min_node = new_node

    def pop(self) -> None:
        pop = self.stack.pop()
        if pop == self.min_node:
            self.min_node = self.min_node.prev

    def top(self) -> int:
        return self.stack[-1].val

    def getMin(self) -> int:
        return self.min_node.val


# Your MinStack object will be instantiated and called as such:
# obj = MinStack()
# obj.push(val)
# obj.pop()
# param_3 = obj.top()
# param_4 = obj.getMin()