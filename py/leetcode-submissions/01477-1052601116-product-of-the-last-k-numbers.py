class ProductOfNumbers:

    def __init__(self):
        self.products, self.zero_indexes = [], []

    def add(self, num: int) -> None:
        if num == 0:
            self.zero_indexes.append(len(self.products))
            num = 1
        new_product = num * (self.products[-1] if self.products else 1)
        self.products.append(new_product)

    def getProduct(self, k: int) -> int:
        first_elem = len(self.products) - k
        if self.zero_indexes and self.zero_indexes[-1] >= first_elem:
            return 0
        if first_elem == 0:
            return self.products[-1]
        return self.products[-1] // self.products[first_elem - 1]


# Your ProductOfNumbers object will be instantiated and called as such:
# obj = ProductOfNumbers()
# obj.add(num)
# param_2 = obj.getProduct(k)