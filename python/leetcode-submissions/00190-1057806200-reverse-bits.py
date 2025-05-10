class Solution:
    def reverseBits(self, n: int) -> int:
        # Initialize the reversed number to 0
        reversed_num = 0
        
        # Iterate over all 32 bits of the given number
        for i in range(32):
            # Left shift the reversed number by 1 and add the last bit of the given number to it
            reversed_num = (reversed_num << 1) | (n & 1)
            # remove the last bit from the original number
            n >>= 1
        
        # Return the reversed number
        return reversed_num
