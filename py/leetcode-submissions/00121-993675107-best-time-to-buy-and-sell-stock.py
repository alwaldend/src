class Solution:
    # Input: prices = [7,1,5,3,6,4]
    # Output: 5
    # Explanation: Buy on day 2 (price = 1) and sell on day 5 (price = 6), profit = 6-1 = 5.
    # Note that buying on day 2 and selling on day 1 is not allowed because you must buy 
    # before you sell.
    def maxProfit(self, prices: List[int]) -> int:
        length = len(prices)
        if length < 2:
            return 0

        left, right, profit = 0, 1, 0
        while right < length:
            current_profit = prices[right] - prices[left]
            is_profitable = current_profit > 0
            if is_profitable and current_profit > profit:
                profit = current_profit
            elif not is_profitable:
                left = right
            
            right += 1

        return profit