class Solution:
    def minSpeedOnTime(self, dist: List[int], hour: float) -> int:
        length = len(dist)
        if hour >= sum(dist):
            return 1
        
        if hour <= length - 1:
            return -1
        
        time_remaining = hour - length + 1
        max_speed = int(max(
            max(dist), dist[-1] // time_remaining + 1
        ))
        min_speed = 1
        result = -1
        while min_speed < max_speed:
            speed = min_speed + (max_speed - min_speed) // 2
            time = dist[-1] / speed + sum(
                (distance + speed - 1) // speed 
                for distance in dist[:-1]
            )
            
            if time > hour:
                min_speed = speed + 1
            else:
                max_speed = speed
        
        return min_speed