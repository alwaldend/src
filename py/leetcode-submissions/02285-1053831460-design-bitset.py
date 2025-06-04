class Bitset:

    def __init__(self, size: int):
        self.bits = [0] * size
        self.ones_count = 0
        self.do_flip = False
        self.size = size

    def fix(self, idx: int) -> None:
        cur = self.bits[idx]
        if self.do_flip and cur == 1:
            self.ones_count -= 1
            self.bits[idx] = 0
        elif not self.do_flip and cur == 0:
            self.ones_count += 1
            self.bits[idx] = 1

    def unfix(self, idx: int) -> None:
        cur = self.bits[idx]
        if self.do_flip and cur == 0:
            self.ones_count += 1
            self.bits[idx] = 1
        elif not self.do_flip and cur == 1:
            self.ones_count -= 1
            self.bits[idx] = 0

    def flip(self) -> None:
        self.do_flip = not self.do_flip

    def all(self) -> bool:
        return self.count() == self.size

    def one(self) -> bool:
        return self.count() > 0

    def count(self) -> int:
        return self.size - self.ones_count if self.do_flip else self.ones_count

    def toString(self) -> str:
        target = (bit ^ 1 for bit in self.bits) if self.do_flip else self.bits
        return "".join(str(num) for num in target)


# Your Bitset object will be instantiated and called as such:
# obj = Bitset(size)
# obj.fix(idx)
# obj.unfix(idx)
# obj.flip()
# param_4 = obj.all()
# param_5 = obj.one()
# param_6 = obj.count()
# param_7 = obj.toString()