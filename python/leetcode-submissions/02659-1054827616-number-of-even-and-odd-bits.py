class Solution:
    def evenOddBit(self, n: int) -> List[int]:
        even, odd = n & 0x55555555, n & 0xaaaaaaaa
        even_count, odd_count = 0, 0
        while even:
            if even & 1 == 1:
                even_count += 1
            even >>= 1
            
        while odd:
            if odd & 1 == 1:
                odd_count += 1
            odd >>= 1
        
        return even_count, odd_count