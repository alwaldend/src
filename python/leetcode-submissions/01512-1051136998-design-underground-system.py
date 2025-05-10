class UndergroundSystem:

    def __init__(self):
        self.from_to = defaultdict(lambda: defaultdict(list))
        self.in_transit = {}

    def checkIn(self, id: int, stationName: str, t: int) -> None:
        self.in_transit[id] = (stationName, t)

    def checkOut(self, id: int, stationName: str, t: int) -> None:
        station_from, time_from = self.in_transit.pop(id)
        self.from_to[station_from][stationName].append(t - time_from)

    def getAverageTime(self, startStation: str, endStation: str) -> float:
        times = self.from_to[startStation][endStation]
        return sum(times) / len(times)


# Your UndergroundSystem object will be instantiated and called as such:
# obj = UndergroundSystem()
# obj.checkIn(id,stationName,t)
# obj.checkOut(id,stationName,t)
# param_3 = obj.getAverageTime(startStation,endStation)