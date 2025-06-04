class Solution:
    def maxSumAfterPartitioning(self, arr: List[int], k: int) -> int:
        min_num = -1
        length = len(arr)
        last_idx = length - 1
        cache = [-1] * length

        def dp(start: int) -> int:
            cached = cache[start]
            if cached != -1:
                return cached
            max_num = min_num
            max_sum = 0
            for i in range(start, min(start + k, length)):
                max_num = max(max_num, arr[i])
                cur_sum = max_num * (i - start + 1)
                if i != last_idx:
                    cur_sum += dp(i + 1)
                max_sum = max(max_sum, cur_sum)
            cache[start] = max_sum
            return max_sum
        
        return dp(0)
            