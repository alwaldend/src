class CustomStack:

    def __init__(self, maxSize: int):
        self._stack = []
        self._max_size = maxSize

    def push(self, x: int) -> None:
        if len(self._stack) != self._max_size:
            self._stack.append(x)

    def pop(self) -> int:
        return self._stack.pop() if self._stack else -1 
        
    def increment(self, k: int, val: int) -> None:
        for i in range(min(k, len(self._stack))):
            self._stack[i] += val


# Your CustomStack object will be instantiated and called as such:
# obj = CustomStack(maxSize)
# obj.push(x)
# param_2 = obj.pop()
# obj.increment(k,val)