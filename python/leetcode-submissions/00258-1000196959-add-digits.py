class Solution:
    def addDigits(self, num: int) -> int:
        sum = num
        while sum > 9:
            current_number, current_sum = sum, 0
            while current_number:
                current_sum += current_number % 10
                current_number //= 10
            sum = current_sum
        return sum