class Solution:
    def permute(self, nums: List[int]) -> List[List[int]]:

        length = len(nums)
        current = []
        current_contains = [False] * length

        def backtrack() -> Generator[None, None, List[int]]:
            if len(current) == length:
                yield tuple(current[:])
                return

            for i in range(length):
                if current_contains[i]:
                    continue

                current_contains[i] = True
                current.append(nums[i])

                yield from backtrack()

                current_contains[i] = False
                current.pop()
            
            return

        return tuple(combination for combination in backtrack()) 