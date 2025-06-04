class Solution:
    def isPalindrome(self, x: int) -> bool:
        if x < 0:
            return False
        
        if x < 10:
            return True
        
        number = []

        while x:
            remainder = x % 10
            x = x // 10
            number.append(remainder)
            
        length = len(number)
        half_index = length // 2

        for i, digit in enumerate(number):
            last_digit = number[length - i - 1]

            if digit != last_digit:
                return False
            
            if i != half_index:
                continue
            
            return True
