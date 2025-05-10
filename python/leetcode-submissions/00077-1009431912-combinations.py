class Solution:
    def combine(self, n: int, k: int) -> List[List[int]]:
        def generate_combinations(elems: List[int], num: int):
            total = len(elems)
            if num > total:
                return
            curr_indices = list(range(num))
            reversed_num = tuple(reversed(range(num)))

            while True:
                yield list(elems[i] for i in curr_indices)
                
                for idx in reversed_num:
                    if curr_indices[idx] != idx + total - num:
                        break
                else:
                    return

                curr_indices[idx] += 1
                for j in range(idx+1, num):
                    curr_indices[j] = curr_indices[j-1] + 1

        return [combination for combination in generate_combinations(tuple(range(1, n+1)), k)]