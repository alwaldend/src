class FrequencyTracker:

    def __init__(self):
        self._num_freq = defaultdict(int)
        self._freq_nums = defaultdict(set)    

    def add(self, number: int) -> None:
        cur_freq = self._num_freq[number]
        self._num_freq[number] = cur_freq + 1
        self._freq_nums[cur_freq].discard(number)
        self._freq_nums[cur_freq + 1].add(number)

    def deleteOne(self, number: int) -> None:
        cur_freq = self._num_freq[number]
        if cur_freq == 0:
            return    
        self._num_freq[number] = cur_freq - 1 
        self._freq_nums[cur_freq].discard(number)
        if cur_freq != 1:
            self._freq_nums[cur_freq - 1].add(number)

    def hasFrequency(self, frequency: int) -> bool:
        return len(self._freq_nums[frequency]) != 0


# Your FrequencyTracker object will be instantiated and called as such:
# obj = FrequencyTracker()
# obj.add(number)
# obj.deleteOne(number)
# param_3 = obj.hasFrequency(frequency)