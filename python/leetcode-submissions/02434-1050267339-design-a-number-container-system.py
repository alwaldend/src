class NumberContainers:

    def __init__(self):
        self.num_indices = defaultdict(list)
        self.num_at_index = {}
        

    def change(self, index: int, number: int) -> None:
        self.num_at_index[index] = number
        heapq.heappush(self.num_indices[number], index)
        

    def find(self, number: int) -> int:
        while self.num_indices[number] and self.num_at_index[self.num_indices[number][0]] != number:
            heapq.heappop(self.num_indices[number])
        
        return self.num_indices[number][0] if len(self.num_indices[number]) > 0 else -1
        


# Your NumberContainers object will be instantiated and called as such:
# obj = NumberContainers()
# obj.change(index,number)
# param_2 = obj.find(number)