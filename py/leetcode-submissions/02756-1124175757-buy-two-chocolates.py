class Solution:
    def buyChoco(self, prices: List[int], money: int) -> int:
        # Assume the Minimum Cost to be Infinity
        min_cost = float('inf')

        # Number of Chocolates
        n = len(prices)

        # Check Every Pair of Chocolates
        for first_choco in range(n):
            for second_choco in range(first_choco + 1, n):
                # Sum of Prices of the Two Chocolates
                cost = prices[first_choco] + prices[second_choco]

                # If the Sum of Prices is Less than the Minimum Cost
                if cost < min_cost:
                    # Update the Minimum Cost
                    min_cost = cost
        
        # We can buy chocolates only if we have enough money
        if min_cost <= money:
            # Return the Amount of Money Left
            return money - min_cost
        else:
            # We cannot buy chocolates. Return the initial amount of money
            return money