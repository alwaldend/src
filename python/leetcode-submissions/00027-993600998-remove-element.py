class Solution:
    # create replace index
    # iterate over nums:
    # - if the current number is equal to val, continue
    # - set nums[replace] to that number, increase the index

    def removeElement(self, nums: List[int], val: int) -> int:
        replace = 0
        for i, number in enumerate(nums):
            if number == val:
                continue
            nums[replace] = number
            replace += 1
        
        return replace