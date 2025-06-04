class Solution:
    def cherryPickup(self, grid: List[List[int]]) -> int:
        n = len(grid)
        m = len(grid[0])

        # Create 3D DP table with initial values of 0
        dp = [[[0] * m for _ in range(m)] for _ in range(n)]

        # Set the starting point value (top-left and top-right corner)
        cherries = 0
        dp[0][0][m - 1] = grid[0][0] + grid[0][m - 1]

        # Iterate through each row from second onwards
        for i in range(1, n):
            # Iterate through each column for robot 1
            for j in range(m):
                # Iterate through each column for robot 2
                for k in range(m):
                    # Skip invalid states:
                    # - Both robots in the same row (j > i)
                    # - Robot 2 left of robot 1 (k < m - i - 1)
                    # - Robot 1 further right than robot 2 (j > k)
                    if j > i or k < m - i - 1 or j > k:
                        continue
                    # Base case: no moves possible, use previous state
                    dp[i][j][k] = dp[i - 1][j][k]
                    # Explore moves for robot 1:
                    # - Up-diagonal with robot 2 at same position
                    if j - 1 >= 0:
                        dp[i][j][k] = max(dp[i][j][k], dp[i - 1][j - 1][k])
                    # - Up-diagonal with robot 2 one step left/right
                    if j - 1 >= 0 and k - 1 >= 0:
                        dp[i][j][k] = max(dp[i][j][k], dp[i - 1][j - 1][k - 1])
                    if j - 1 >= 0 and k + 1 < m:
                        dp[i][j][k] = max(dp[i][j][k], dp[i - 1][j - 1][k + 1])
                    # Explore moves for robot 2:
                    # - Up-diagonal with robot 1 at same position
                    if j + 1 < m:
                        dp[i][j][k] = max(dp[i][j][k], dp[i - 1][j + 1][k])
                    # - Up-diagonal with robot 1 one step left/right
                    if j + 1 < m and k - 1 >= 0:
                        dp[i][j][k] = max(dp[i][j][k], dp[i - 1][j + 1][k - 1])
                    if j + 1 < m and k + 1 < m:
                        dp[i][j][k] = max(dp[i][j][k], dp[i - 1][j + 1][k + 1])
                    # Explore horizontal moves for both robots:
                    # - Both robots move left
                    if k - 1 >= 0:
                        dp[i][j][k] = max(dp[i][j][k], dp[i - 1][j][k - 1])
                    # - Both robots move right
                    if k + 1 < m:
                        dp[i][j][k] = max(dp[i][j][k], dp[i - 1][j][k + 1])
                    # Add cherries only if robots are in different positions
                    if j != k:
                        dp[i][j][k] += grid[i][j] + grid[i][k]
                    else:
                        dp[i][j][k] += grid[i][j]  # Only one robot picks if they land in the same cell
                    # Update maximum cherries collected so far
                    cherries = max(cherries, dp[i][j][k])

        return cherries
