class Solution:
    def addBinary(self, a: str, b: str) -> str:
        answer = []
        carry = 0
        for char1, char2 in zip_longest(reversed(a), reversed(b)):
            num1, num2 = int(char1) if char1 else 0, int(char2) if char2 else 0
            cur_sum = num1 + num2 + carry
            carry = 0
            if cur_sum > 1:
                cur_sum -= 2
                carry = 1
            
            answer.append(str(cur_sum))
        
        if carry:
            answer.append("1")

        return "".join(reversed(answer))