class SnapshotArray:

    def __init__(self, length: int):
        self._cur_snap = 0
        self._elems = [[(0, 0)] for _ in range(length)]

    def set(self, index: int, val: int) -> None:
        values = self._elems[index]
        if values[-1][0] == self._cur_snap:
            values.pop()
        values.append((self._cur_snap, val))

    def snap(self) -> int:
        self._cur_snap += 1
        return self._cur_snap - 1 

    def get(self, index: int, snap_id: int) -> int:
        for cur_snap_id, val in reversed(self._elems[index]):
            if cur_snap_id > snap_id:
                continue
            return val
        
        return -1
        


# Your SnapshotArray object will be instantiated and called as such:
# obj = SnapshotArray(length)
# obj.set(index,val)
# param_2 = obj.snap()
# param_3 = obj.get(index,snap_id)