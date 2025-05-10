class Solution:
    def maxConsecutiveAnswers(self, answerKey: str, k: int) -> int:
        length = len(answerKey)
        char_t, char_f = "T", "F"
        max_size = k
        count = defaultdict(int)
        
        for char in answerKey[:k]:
            count[char] += 1

        left = 0
        for right in range(k, length):
            count[answerKey[right]] += 1
            
            while min(count[char_t], count[char_f]) > k: 
                count[answerKey[left]] -= 1
                left += 1
            
            size = right - left + 1
            if size > max_size:
                max_size = size
                    
        return max_size



        
