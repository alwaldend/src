class Solution:
    # non-decreasing order, so to remove the duplicates we just need to remove all 
    #    consequent duplicates
    # create replace index, set it to 1 - the first element is always unique
    # check if length is more than 1 to avoid out-of-bounds - 
    #   if the length is one, just return 1
    # iterate over nums starting from the second element:
    # - if the current number is not equal to the previous, 
    #   set nums[replace] to it, move the replace index
    # - if the current number is equal to the previous one, continue 
    # return replace 
    def removeDuplicates(self, nums: List[int]) -> int:
        replace = 1
        for i, number in enumerate(nums[1:], 1): 
            if number == nums[i-1]:
                continue
            nums[replace] = number
            replace += 1
        
        return replace