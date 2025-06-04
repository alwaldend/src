class Solution:
    # non-decreasing order, so to remove the duplicates we just need to remove all 
    #    consequent duplicates
    # create replace index, set it to 0
    # create a set of duplicates
    # iterate over nums:
    # - if the number is in the set, continue
    # - if the number is not in the set, set nums[replace] to that number, add it to the set
    # return replace
    def removeDuplicates(self, nums: List[int]) -> int:
        replace, duplicates = 0, set()
        for i, number in enumerate(nums): 
            if number in duplicates:
                continue
            nums[replace] = number
            replace += 1
            duplicates.add(number)
        
        return replace