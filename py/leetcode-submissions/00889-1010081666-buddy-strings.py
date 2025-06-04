class Solution:
    def buddyStrings(self, s: str, goal: str) -> bool:
        length_1, length_2 = len(s), len(goal)

        if length_1 != length_2:
            return False
        
        if s == goal:
            freq = defaultdict(int)
            for char in s:
                freq[char] += 1
                if freq[char] == 2:
                    return True
            return False
        
        swap_1, swap_2 = -1, -1
        
        for i in range(length_1):
            char_1, char_2 = s[i], goal[i]
            
            if char_1 == char_2:
                continue
            
            if swap_1 == -1:
                swap_1 = i
            elif swap_2 == -1:
                swap_2 = i
            else:
                return False
            

        return swap_2 != -1 and s[swap_1] == goal[swap_2] and s[swap_2] == goal[swap_1]