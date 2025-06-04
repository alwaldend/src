class Solution:
    def isPalindrome(self, x: int) -> bool:
        if x < 0:
            return False
        
        if x < 10:
            return True
        
        number = []

        while True:
            remainder = x % 10
            x = int(x / 10)
            number.append(remainder)
            if not x:
                break
            
        
        number.reverse()
        print(number)
        length = len(number)
        half_index = int(length / 2)

        for i, digit in enumerate(number):
            last_digit = number[length - i - 1]

            if digit != last_digit:
                return False
            
            if i != half_index:
                continue
            
            return True
