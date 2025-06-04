class Solution:
    # n = 2, batteries = [3,3,3], Output: 4
    # n = 2, batteries = [1,1,1,1], Output: 2
    def maxRunTime(self, n: int, batteries: List[int]) -> int:
        length = len(batteries)

        if length < n:
            return 0
        
        if length == n:
            return min(batteries)

        batteries.sort()
        extra = sum(batteries[:-n])
        live = batteries[-n:]

        
        # We increase the total running time using 'extra' by increasing 
        # the running time of the computer with the smallest battery.
        for i in range(n - 1):
            # If the target running time is between live[i] and live[i + 1].
            if extra // (i + 1) < live[i + 1] - live[i]:
                return live[i] + extra // (i + 1)
            
            # Reduce 'extra' by the total power used.
            extra -= (i + 1) * (live[i + 1] - live[i])
        
        # If there is power left, we can increase the running time 
        # of all computers.
        return live[-1] + extra // n