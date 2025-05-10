class Solution:
    def PredictTheWinner(self, nums: List[int]) -> bool:
        length = len(nums)
        is_even = length % 2 == 0

        if length < 3:
            return True
        
        @cache
        def max_diff(left: int, right: int) -> int:
            left_num, right_num = nums[left], nums[right]
            
            if left == right:
                return left_num
            
            score_by_left = left_num - max_diff(left + 1, right)
            score_by_right = right_num - max_diff(left, right - 1)
            return max(score_by_left, score_by_right)

        return max_diff(0, length - 1) >= 0
