class Cashier:

    def __init__(self, n: int, discount: int, products: List[int], prices: List[int]):
        self.freq, self.discount = n, (100 - discount) / 100
        self.cur = 0
        self.products = {id: price for id, price in zip(products, prices)}

    def getBill(self, product: List[int], amount: List[int]) -> float:
        self.cur += 1
        multiplier = 1
        if self.cur == self.freq:
            self.cur, multiplier = 0, self.discount
        total = sum(amount[i] * self.products[product[i]] for i in range(len(product)))
        return total * multiplier
        


# Your Cashier object will be instantiated and called as such:
# obj = Cashier(n, discount, products, prices)
# param_1 = obj.getBill(product,amount)