class Solution:
    def findLeastNumOfUniqueInts(self, arr: List[int], k: int) -> int:
        counter = defaultdict(int)
        for num in arr:
            counter[num] += 1
        nums = [(count, num) for num, count in counter.items()]
        nums.sort()
        removed = 0
        for count, num in nums:
            if k < count:
                break
            k -= count
            removed += 1
        return len(counter) - removed