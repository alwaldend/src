class ProductOfNumbers:

    def __init__(self):
        self.products, self.last_zero, self.length = [], -1, 0

    def add(self, num: int) -> None:
        if num == 0:
            self.last_zero = self.length
            num = 1
        self.products.append(num * (self.products[-1] if self.products else 1))
        self.length += 1

    def getProduct(self, k: int) -> int:
        first_elem = self.length - k
        if self.last_zero >= first_elem:
            return 0
        if first_elem == 0:
            return self.products[-1]
        return self.products[-1] // self.products[first_elem - 1]


# Your ProductOfNumbers object will be instantiated and called as such:
# obj = ProductOfNumbers()
# obj.add(num)
# param_2 = obj.getProduct(k)