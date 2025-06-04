class Solution:
    def coinChange(self, coins: List[int], amount: int) -> int:
        dp = [amount + 1] * (amount + 1)
        dp[0] = 0

        for current_amount in range(1, amount + 1):
            for coin in coins:
                diff = current_amount - coin
                if diff < 0:
                    continue
                
                dp[current_amount] = min(dp[current_amount], dp[diff] + 1)

        return dp[-1] if dp[-1] != amount + 1 else -1
