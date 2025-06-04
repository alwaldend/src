class Solution:
    def coinChange(self, coins: List[int], amount: int) -> int:
        if amount == 0:
            return 0
        
        @cache
        def dp(coins_value: int) -> int:
            if coins_value == amount:
                return 0
            
            if coins_value > amount:
                return -1
            
            min_coins_count = -1

            for coin in coins:
                new_count = 1 + dp(coin + coins_value)
                if new_count == 0:
                    continue
                
                if min_coins_count == -1 or new_count < min_coins_count:
                    min_coins_count = new_count
            
            return min_coins_count

        return dp(0)
