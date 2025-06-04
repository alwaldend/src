class Solution:
    def sortArrayByParity(self, nums: List[int]) -> List[int]:
        length = len(nums)
        left, right = 0, length - 1
        answer = [None] * length
        for num in nums:
            if num % 2 == 0:
                answer[left] = num
                left += 1
            else:
                answer[right] = num
                right -= 1
        return answer