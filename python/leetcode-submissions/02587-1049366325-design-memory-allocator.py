class Allocator:

    def __init__(self, n: int):
      self._units = [1] * n
      self._units[0] = n
      self._id_units = defaultdict(list)
      self._units_count = n

    def find_avail_units(self, start: int, size: int) -> Tuple[int, int]:
        i = start
        count = 0
        while i < self._units_count and i - start < size:
            units_avail = self._units[i]
            if units_avail < 0:
                return i + units_avail if i == start else i, count
            i += units_avail
            count += units_avail

        return i, count
            

    def allocate(self, size: int, mID: int) -> int:
        i = 0
        units_start, units_count = None, 0
        while i < self._units_count and units_count < size:
            units_avail = self._units[i]
            if units_avail < 0:
                i += -units_avail
                units_start, units_count = None, 0
                continue

            if units_start is None:
                units_start = i
            units_count += units_avail
            i += units_avail

        if units_count < size:
            return -1 

        self._units[units_start] = -size
        if units_count > size:
            self._units[units_start + size] = units_count - size
        self._id_units[mID].append(units_start)
        return units_start

    def free(self, mID: int) -> int:
        count = 0
        conseq_units = self._id_units[mID]
        while conseq_units:
            units_start = conseq_units.pop()
            units_freed = -self._units[units_start]
            count += units_freed
            self._units[units_start] = units_freed

        return count


# Your Allocator object will be instantiated and called as such:
# obj = Allocator(n)
# param_1 = obj.allocate(size,mID)
# param_2 = obj.free(mID)