class Solution:
    def letterCombinations(self, digits: str) -> List[str]:
        length = len(digits)
        if length == 0:
            return []
        
        digit_map: Dict[int, str] = {
            1: [],
            2: ["a", "b", "c"],
            3: ["d", "e", "f"],
            4: ["g", "h", "i"],
            5: ["j", "k", "l"],
            6: ["m", "n", "o"],
            7: ["p", "q", "r", "s"],
            8: ["t", "u", "v"],
            9: ["w", "x", "y", "z"],
            0: [" "]
        }
        result: List[List[str]] = [
            digit_map[int(digit)] for digit in digits
        ]
        return ["".join(i) for i in product(*result)]