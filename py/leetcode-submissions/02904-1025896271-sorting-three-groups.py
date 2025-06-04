class Solution:
    def minimumOperations(self, nums: List[int]) -> int:
        nums_count = len(nums)
        
        if nums_count == 1:
            return 0
            
        @cache
        def dp(i: int, start_group: int) -> int:
            if i == nums_count:
                return 0
            
            curr_group = nums[i] - 1
            actions_min = None
            
            for group_available in range(start_group, 3):
                actions = dp(i + 1, group_available) + (0 if group_available == curr_group else 1)
                if actions_min is None or actions < actions_min:
                    actions_min = actions
            
            return actions_min
        
        return dp(0, 0)