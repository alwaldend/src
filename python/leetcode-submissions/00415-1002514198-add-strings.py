class Solution:
    def addStrings(self, num1: str, num2: str) -> str:
        result = []
        index_1, index_2 = len(num1) - 1, len(num2) - 1


        carry = 0
        while index_1 >= 0 or index_2 >= 0 or carry:
            digit_1 = num1[index_1] if index_1 >= 0 else 0
            digit_2 = num2[index_2] if index_2 >= 0 else 0
            digit = int(digit_1) + int(digit_2) + carry
            if digit > 9:
                carry = 1
                digit %= 10
            else:
                carry = 0

            result.append(str(digit))
            index_1 -= 1
            index_2 -= 1 

        return "".join(reversed(result))