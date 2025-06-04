class MyHashMap:
    def __init__(self):
        self.data = [None] * 1000001

    def put(self, key: int, val: int) -> None:
        self.data[key] = val
        
    def get(self, key: int) -> int:
        val = self.data[key]
        return -1 if val is None else val

    def remove(self, key: int) -> None:
        self.data[key] = None
        


# Your MyHashMap object will be instantiated and called as such:
# obj = MyHashMap()
# obj.put(key,value)
# param_2 = obj.get(key)
# obj.remove(key)