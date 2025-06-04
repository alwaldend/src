class Solution:
    def romanToInt(self, s: str) -> int:
        result = 0
        skip = False
        for i, symbol in enumerate(s):
            if skip:
                skip = False
                continue

            next = ""
            if i < len(s) - 1:
                next = s[i+1]
            
            if symbol == "I" and next == "V":
                result += 4
                skip = True
                continue

            if symbol == "I" and next == "X":
                result += 9
                skip = True
                continue
            
            if symbol == "I":
                result += 1
                continue

            if symbol == "X" and next == "L":
                result += 40
                skip = True
                continue

            if symbol == "X" and next == "C":
                result += 90
                skip = True
                continue
            
            if symbol == "X":
                result += 10
                continue

            if symbol == "C" and next == "D":
                result += 400
                skip = True
                continue
            
            if symbol == "C" and next == "M":
                result += 900
                skip = True
                continue
            
            if symbol == "C":
                result += 100
                continue
            
            if symbol == "V":
                result += 5
                continue
            
            if symbol == "D":
                result += 500
                continue
            
            if symbol == "M":
                result += 1000
                continue
            
            if symbol == "L":
                result += 50
                continue
            
            raise Exception(f"unexpected situation: {symbol}, {next}")

        return result
            