class Solution:
    def rearrangeArray(self, nums: List[int]) -> List[int]:
        res = []
        pos, neg = [], []
        length = len(nums)
        for i in range(length):
            num = nums[i]
            if num > 0:
                pos.append(num)
            else:
                neg.append(num)
        pos.reverse()
        neg.reverse()
        while pos and neg:
            res.extend((pos.pop(), neg.pop()))
        return res