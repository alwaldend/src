class Solution:
    def countCompleteSubarrays(self, nums: List[int]) -> int:
        length = len(nums)
        elems_count = len(set(nums))
        
        if elems_count == length:
            return 1
        
        if elems_count == 1:
            return length + sum(i for i in range(1, length))
        
        result, elems, min_j = 0, defaultdict(int), 0
        for i in range(length):
            left = nums[i]
            
            for j in range(min_j, length):
                right = nums[j]
                elems[right] += 1
                
                if len(elems) != elems_count:
                    continue
                
                if elems[right] == 1:
                    elems.pop(right)
                else:
                    elems[right] -= 1 
    
                result += length - j
                min_j = j
                break
            else:
                return result
            
            if elems[left] == 1:
                elems.pop(left)
            else:
                elems[left] -= 1
                    
        return result