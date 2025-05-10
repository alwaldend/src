class Solution:
    def longestCommonPrefix(self, strs: List[str]) -> str:
        answer = []
        strings = sorted(strs)
        first = strings[0]
        last = strings[-1]
        min_length = min(len(first), len(last))
        for i in range(min_length):
            first_symbol = first[i]
            last_symbol = last[i]
            if first_symbol == last_symbol:
                answer.append(first_symbol)
                continue
            
            return "".join(answer)

        return "".join(answer)


        