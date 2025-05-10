class Solution:
    def frequencySort(self, s: str) -> str:
        counter = defaultdict(int)
        for char in s:
            counter[char] += 1
        pq = [(-freq, char) for char, freq in counter.items()]
        heapq.heapify(pq)
        result = []
        while pq:
            freq, char = heapq.heappop(pq)
            result.append(char * -freq)
        return "".join(result)