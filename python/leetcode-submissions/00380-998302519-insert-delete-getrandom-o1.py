class RandomizedSet:

    def __init__(self):
        self._set = set()
        self._items = []
        self._indexes = {}
        

    def insert(self, val: int) -> bool:
        is_in = val in self._set
        if not is_in:
            self._set.add(val)
            self._items.append(val)
            self._indexes[val] = len(self._items) - 1
        return not is_in

    def remove(self, val: int) -> bool:
        if val not in self._set:
            return False

        last = self._items[-1]
        val_index = self._indexes[val]
        self._items[val_index] = last
        self._indexes[last] = val_index
        
        self._set.remove(val)
        self._items.pop()
        self._indexes.pop(val)

        return True

    def getRandom(self) -> int:
        return random.choice(self._items)


# Your RandomizedSet object will be instantiated and called as such:
# obj = RandomizedSet()
# param_1 = obj.insert(val)
# param_2 = obj.remove(val)
# param_3 = obj.getRandom()