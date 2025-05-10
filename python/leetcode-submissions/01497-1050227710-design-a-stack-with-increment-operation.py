class CustomStack:

    def __init__(self, maxSize: int):
        self._stack = []
        self._max_size = maxSize
        self._size = 0

    def push(self, x: int) -> None:
        if self._size == self._max_size:
            return
        self._stack.append(x)
        self._size += 1

    def pop(self) -> int:
        if self._size == 0:
            return -1
        last = self._stack.pop()
        self._size -= 1
        return last

    def increment(self, k: int, val: int) -> None:
        for i in range(min(k, self._size)):
            self._stack[i] += val


# Your CustomStack object will be instantiated and called as such:
# obj = CustomStack(maxSize)
# obj.push(x)
# param_2 = obj.pop()
# obj.increment(k,val)