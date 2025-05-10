class Solution:
    def countPairs(self, nums: List[int], target: int) -> int:
        pairs_count = 0
        nums_count = len(nums)
        for i in range(nums_count): 
            num1 = nums[i]
            for j in range(i + 1, nums_count):
                if num1 + nums[j] < target:
                    pairs_count += 1
        
        return pairs_count