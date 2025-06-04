class Solution:
    def asteroidCollision(self, asteroids: List[int]) -> List[int]:
        length = len(asteroids)

        if length < 2:
            return asteroids

        stack = [asteroids[0]]
        for i in range(1, length):
            asteroid = asteroids[i]
            asteroid_last = stack[-1] if stack else 0
            stack.append(asteroid)

            if asteroid > 0 or (
                asteroid < 0 and asteroid_last < 0
            ):
                continue

            while len(stack) >= 2 and stack[-1] < 0 and stack[-2] > 0:
                last, prev = abs(stack[-1]), abs(stack[-2])
                if last == prev:
                    stack.pop()
                elif last > prev:
                    stack[-1], stack[-2] = stack[-2], stack[-1]

                stack.pop()

        return stack