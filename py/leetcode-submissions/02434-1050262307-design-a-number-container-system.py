import sortedcontainers

class NumberContainers:

    def __init__(self):
        self._idx_number = {}
        self._number_idx = defaultdict(sortedcontainers.SortedSet)

    def change(self, index: int, number: int) -> None:
        cur_number = self._idx_number.get(index, -1)
        if cur_number != -1:
            self._number_idx[cur_number].remove(index)
        self._idx_number[index] = number
        self._number_idx[number].add(index)

    def find(self, number: int) -> int:
        ids = self._number_idx[number]
        return ids[0] if ids else -1


# Your NumberContainers object will be instantiated and called as such:
# obj = NumberContainers()
# obj.change(index,number)
# param_2 = obj.find(number)