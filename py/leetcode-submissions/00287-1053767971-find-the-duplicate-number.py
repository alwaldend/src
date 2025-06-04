class Solution:
    def findDuplicate(self, nums: List[int]) -> int:
        freqs = [False] * (10**5 + 1)
        for num in nums:
            if freqs[num] == True:
                return num
            freqs[num] = True

        raise Exception()