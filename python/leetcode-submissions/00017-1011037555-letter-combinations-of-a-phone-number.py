class Solution:
    def letterCombinations(self, digits: str) -> List[str]:
        length = len(digits)
        if length == 0:
            return []

        digit_map: Dict[int, str] = {
            "1": [],
            "2": ["a", "b", "c"],
            "3": ["d", "e", "f"],
            "4": ["g", "h", "i"],
            "5": ["j", "k", "l"],
            "6": ["m", "n", "o"],
            "7": ["p", "q", "r", "s"],
            "8": ["t", "u", "v"],
            "9": ["w", "x", "y", "z"],
            "0": [" "]
        }

        current = []
        
        def backtrack(digit: int) -> Generator[None, None, str]:
            if digit == length:
                yield "".join(current[:])
                return
            
            for char in digit_map[digits[digit]]:
                current.append(char)
                yield from backtrack(digit + 1)
                current.pop()

        return tuple(combination for combination in backtrack(0))