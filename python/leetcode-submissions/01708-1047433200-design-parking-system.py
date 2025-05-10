class ParkingSystem:

    def __init__(self, big: int, medium: int, small: int):
        self._slots = [0, big, medium, small]

    def addCar(self, carType: int) -> bool:
        slots_avail = self._slots[carType]
        if not slots_avail:
            return False
        self._slots[carType] = slots_avail - 1
        return True


# Your ParkingSystem object will be instantiated and called as such:
# obj = ParkingSystem(big, medium, small)
# param_1 = obj.addCar(carType)