class Solution:
    # create two indexes, current and replace
    # create non_val_count
    # while to-be-replaced has not reached the end:
    # - if the current number is a regular number, move the current index, 
    #   increase non_val_count
    # - if the replace index is equal or less than the current, move it and continue
    # - if the replace number is a non-regular number, move the replace index and continue
    # - if the current number is a non-regular, replace it with the replace number, 
    #   replace the replace number with val, move both indexes
    # return the current index + 1

    def removeElement(self, nums: List[int], val: int) -> int:
        current, replace, val_count, length = 0, 0, 0, len(nums)
        non_val_count = 0
        
        while replace < length and current < length:
            current_number, replace_number = nums[current], nums[replace]
            
            if current_number != val:
                current += 1
                non_val_count += 1
                continue

            if replace <= current:
                replace = current + 1
                continue

            if replace_number == val:
                replace += 1
                continue

            nums[current], nums[replace] = replace_number, val
            replace += 1
            current += 1
            non_val_count += 1
            
        return non_val_count