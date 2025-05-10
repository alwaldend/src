class Solution:
    def minTaps(self, n: int, ranges: List[int]) -> int:
        inf = float("inf")
        dp = [inf] * (n + 1)
        dp[0] = 0
        
        for i in range(n + 1):
            cur_range = ranges[i]
            tap_start, tap_end = max(0, i - cur_range), min(n, i + cur_range)
            
            for j in range(tap_start, tap_end + 1):
                dp[tap_end] = min(dp[tap_end], dp[j] + 1)
        
        min_taps = dp[n]
        return -1 if min_taps == inf else min_taps 