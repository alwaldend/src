class SmallestInfiniteSet:

    def __init__(self):
        self.heap, self.next_num, self.nums = [], 1, set() 

    def popSmallest(self) -> int:
        if self.heap:
            val = heappop(self.heap)
            self.nums.remove(val)
            return val
        self.next_num += 1
        return self.next_num - 1

    def addBack(self, num: int) -> None:
        if num in self.nums or num >= self.next_num:
            return
        self.nums.add(num)
        heappush(self.heap, num)


# Your SmallestInfiniteSet object will be instantiated and called as such:
# obj = SmallestInfiniteSet()
# param_1 = obj.popSmallest()
# obj.addBack(num)