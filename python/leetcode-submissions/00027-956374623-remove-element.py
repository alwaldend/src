class Solution:
    def removeElement(self, nums: List[int], val: int) -> int:
        length = len(nums)

        if not length: 
            return 0

        first = nums[0]
        if length == 1 and first == val:
            nums.pop()
            return 0
        if length == 1 and first != val:
            return 1
        
        not_equal_index = -1 if first == val else 0
        for i, number in enumerate(nums[1:], 1):
            if number == val:
                continue

            not_equal_index += 1
            nums[not_equal_index] = number

        return not_equal_index + 1
            

