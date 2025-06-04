class TimeMap:

    def __init__(self):
        self._cache = defaultdict(list)
        

    def set(self, key: str, value: str, timestamp: int) -> None:
        self._cache[key].append((timestamp, value))

    def get(self, key: str, timestamp: int) -> str:
        for cur_timestamp, value in reversed(self._cache[key]):
            if cur_timestamp > timestamp:
                continue
            return value

        return ""

# Your TimeMap object will be instantiated and called as such:
# obj = TimeMap()
# obj.set(key,value,timestamp)
# param_2 = obj.get(key,timestamp)