class Solution:
    def binaryGap(self, n: int) -> int:
        i, first_bit = 0, None
        max_distance = 0
        while n:
            is_one = n & 1 == 1
            if is_one and first_bit is None:
                first_bit = i
            elif is_one:
                max_distance, first_bit = max(max_distance, i - first_bit), i
            n >>= 1
            i += 1
        
        return max_distance 