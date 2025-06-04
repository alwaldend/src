class Solution:
    # [1,1,3,4], [1,2,3,4] -> [2]
    def findDisappearedNumbers(self, nums: List[int]) -> List[int]:
        nums.sort()
        result = []
        length = len(nums)
        j = 0
        for i in range(1, length + 1):
            while j < length and nums[j] < i:
                j += 1
            
            if j < length and nums[j] == i:
                continue

            result.append(i)

        return result
