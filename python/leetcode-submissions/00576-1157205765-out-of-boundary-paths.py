class Solution:

    def findPaths(self, m: int, n: int, N: int, x: int, y: int) -> int:

        M = 1000000000 + 7

        dp = [[0] * n for _ in range(m)]

        dp[x][y] = 1

        count = 0



        for moves in range(1, N + 1):

            temp = [[0] * n for _ in range(m)]



            for i in range(m):

                for j in range(n):

                    if i == m - 1:

                        count = (count + dp[i][j]) % M

                    if j == n - 1:

                        count = (count + dp[i][j]) % M

                    if i == 0:

                        count = (count + dp[i][j]) % M

                    if j == 0:

                        count = (count + dp[i][j]) % M

                    temp[i][j] = (

                        ((dp[i - 1][j] if i > 0 else 0) + (dp[i + 1][j] if i < m - 1 else 0)) % M +

                        ((dp[i][j - 1] if j > 0 else 0) + (dp[i][j + 1] if j < n - 1 else 0)) % M

                    ) % M



            dp = temp



        return count





