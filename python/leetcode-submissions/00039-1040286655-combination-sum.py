class Solution:
    def combinationSum(self, candidates: List[int], target: int) -> List[List[int]]:
        result = set()
        cur_nums = []

        def backtrack(cur_sum: int) -> None:
            if cur_sum == target:
                result.add(tuple(sorted(cur_nums[:])))
            if cur_sum >= target:
                return

            for num in candidates:
                cur_nums.append(num)
                backtrack(num + cur_sum)
                cur_nums.pop() 

        backtrack(0)

        return result