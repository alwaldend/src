class Solution:
    def minBitFlips(self, start: int, goal: int) -> int:
        flips_count = 0
        while start or goal:
            start_is_one, goal_is_one = start & 1, goal & 1
            if start_is_one and not goal_is_one or (not start_is_one and goal_is_one):
                flips_count += 1
            start >>= 1
            goal >>= 1
        return flips_count