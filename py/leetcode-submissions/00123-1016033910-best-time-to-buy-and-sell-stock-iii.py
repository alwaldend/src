class Solution:
	def maxProfit(self, prices: List[int]) -> int:
				
		'''
		dp_2_hold: max profit with 2 transactions, and in hold state
		dp_2_not_hold: max profit with 2 transactions, and not in hold state
		
		dp_1_hold: max profit with 1 transaction, and in hold state
		dp_1_not_hold: max profit with 1 transaction, and not in hold state
		
		Note: it is impossible to have stock in hand and sell on first day, therefore -infinity is set as initial profit value for hold state
		'''
		
		dp_2_hold, dp_2_not_hold = -float('inf'), 0
		dp_1_hold, dp_1_not_hold = -float('inf'), 0
		
		for stock_price in prices:
				
			# either keep being in not-hold state, or sell with stock price today
			dp_2_not_hold = max( dp_2_not_hold, dp_2_hold + stock_price )
	
			# either keep being in hold state, or just buy with stock price today ( add one more transaction )
			dp_2_hold = max( dp_2_hold, dp_1_not_hold - stock_price )
				
			# either keep being in not-hold state, or sell with stock price today
			dp_1_not_hold = max( dp_1_not_hold, dp_1_hold + stock_price )
	
			# either keep being in hold state, or just buy with stock price today ( add one more transaction )
			dp_1_hold = max( dp_1_hold, 0 - stock_price )
				
		return dp_2_not_hold