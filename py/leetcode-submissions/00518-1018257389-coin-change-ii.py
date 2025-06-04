class Solution:
    def change(self, amount: int, coins: List[int]) -> int:
        coins_count = len(coins)
        memo = [[-1] * (amount + 1) for _ in range(coins_count)]
        
        def dp(i: int, amount: int) -> int:
            if amount == 0:
                return 1
            if i == coins_count:
                return 0
            if memo[i][amount] != -1:
                return memo[i][amount]

            value = None

            if coins[i] > amount:
                value = dp(i + 1, amount)
            else:
                value = dp(i, amount - coins[i]) + dp(i + 1, amount)

            memo[i][amount] = value
            return value

        return dp(0, amount)