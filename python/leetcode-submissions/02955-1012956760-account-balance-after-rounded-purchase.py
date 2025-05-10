class Solution:
    def accountBalanceAfterPurchase(self, purchaseAmount: int) -> int:
        remainder = purchaseAmount % 10
        if remainder == 0:
            return 100 - purchaseAmount

        if remainder >= 5:
            return 100 - ((purchaseAmount // 10) + 1) * 10 
        
        return 100 - (purchaseAmount // 10) * 10