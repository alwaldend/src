class Solution:
    def intToRoman(self, num: int) -> str:
        result: List[str] = []
        stack: List[int] = deque([
            1, 4, 5, 9, 10, 40, 50, 90, 100, 400, 500, 900, 1000
        ])
        values = {
            1000: "M",
            900: "CM",
            500: "D",
            400: "CD",
            100: "C",
            90: "XC",
            50: "L",
            40: "XL",
            10: "X",
            9: "IX",
            5: "V",
            4: "IV",  
            1: "I"
        }

        while num:
            roman = stack[-1]
            if num < roman:
                stack.pop()
                continue

            result.append(values[roman])
            num -= roman

        return "".join(result)