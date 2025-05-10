class OrderedStream:

    def __init__(self, n: int):
        self.values = [None] * n
        self.length = n
        self.cur = 0

    def insert(self, idKey: int, value: str) -> List[str]:
        self.values[idKey - 1] = value
        for i in range(self.cur, self.length):
            if self.values[i] is not None:
                continue
            
            self.cur, answer = i, self.values[self.cur:i]
            return answer

        return self.values[self.cur:]



# Your OrderedStream object will be instantiated and called as such:
# obj = OrderedStream(n)
# param_1 = obj.insert(idKey,value)