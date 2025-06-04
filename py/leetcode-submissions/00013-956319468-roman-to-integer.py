class Solution:
    def romanToInt(self, s: str) -> int:
        result = 0
        skip = False

        values = {
            "I": 1,
            "V": 5,
            "X": 10,
            "L": 50,
            "C": 100,
            "D": 500,
            "M": 1000
        }
        special_cases = {
            "I": ["V", "X"],
            "X": ["L", "C"],
            "C": ["D", "M"]
        }

        for i, symbol in enumerate(s):
            if skip:
                skip = False
                continue

            next = None
            if i < (len(s) - 1):
                next = s[i+1]

            if next in special_cases.get(symbol, []):
                skip = True
                result += values[next] - values[symbol]
                continue
            
            result += values[symbol]

        return result
            