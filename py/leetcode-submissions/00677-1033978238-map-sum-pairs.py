class MapSum:

    def __init__(self):
        self._map = {}

    def insert(self, key: str, val: int) -> None:
        self._map[key] = val

    def sum(self, prefix: str) -> int:
        return sum(value for key, value in self._map.items() if key.startswith(prefix))


# Your MapSum object will be instantiated and called as such:
# obj = MapSum()
# obj.insert(key,val)
# param_2 = obj.sum(prefix)