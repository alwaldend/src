class Solution:
    def xorOperation(self, n: int, start: int) -> int:
        return reduce(lambda total, i: total ^ (start + 2 * i), chain((0, ), range(n)))