class Solution:
    def findLongestChain(self, pairs: List[List[int]]) -> int:
        pairs_count = len(pairs)
        if pairs_count < 2:
            return pairs_count
        
        pairs.sort()
        
        @cache
        def dp(curr_pair: int) -> int:
            right = pairs[curr_pair][1]
            max_length = 1

            for new_pair in range(curr_pair + 1, pairs_count):
                new_left = pairs[new_pair][0]
                new_length = dp(new_pair) + (1 if new_left > right else 0)

                if new_length > max_length:
                    max_length = new_length
            
            return max_length
        
        return dp(0)