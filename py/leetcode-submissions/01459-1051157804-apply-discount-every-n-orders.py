class Cashier:

    def __init__(self, n: int, discount: int, products: List[int], prices: List[int]):
        self.freq, self.discount = n, (100 - discount) / 100
        self.cur = 0
        self.products = {id: price for id, price in zip(products, prices)}

    def getBill(self, product: List[int], amount: List[int]) -> float:
        self.cur += 1
        total = sum(amount[i] * self.products[product[i]] for i in range(len(product)))
        if self.cur < self.freq:
            return total
        self.cur = 0
        return total * self.discount
        


# Your Cashier object will be instantiated and called as such:
# obj = Cashier(n, discount, products, prices)
# param_1 = obj.getBill(product,amount)