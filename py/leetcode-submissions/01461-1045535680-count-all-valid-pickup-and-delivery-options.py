class Solution:
    def countOrders(self, n: int) -> int:
        # P1, D1                 [01] -> (P1, D1)
        # P1, D1, P2, D2         [06] -> (P1, P2, D1, D2), (P1, P2, D2, D1), 
        #                                (P1, D1, P2, D2), (P2, P1, D1, D2), 
        #                                (P2, P1, D2, D1), (P2, D2, P1, D1).
        # P1, D1, P2, D2, P3, D3 [90]
        
        ways_count = 1
        mod = 10**9 + 7
        for order in range(2, n + 1):
            pos_avail_count = 1 + (order - 1) * 2
            cur_ways_count = pos_avail_count * (pos_avail_count + 1) // 2
            ways_count = (ways_count * cur_ways_count) % mod

        return ways_count