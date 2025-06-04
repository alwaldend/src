class Solution:
    # have three indexes: nums1 (end of array 1), nums2 (from the end), 
    # and current (nums1 from the end) 
    # if nums2 number is bigger or equal than nums1 number, put the number at the current index,
    # move both indexes
    # if nums1 number is smaller than nums2 number, put the number at the current index, 
    # move both indexes  
    def merge(self, nums1: List[int], m: int, nums2: List[int], n: int) -> None:
        """
        Do not return anything, modify nums1 in-place instead.
        """

        nums1_i, nums2_i = m - 1, n - 1
        for i in range(m + n - 1, -1, -1):
            nums1_number  = nums1[nums1_i] if nums1_i >= 0 else nums2[0] - 1
            nums2_number = nums2[nums2_i] if nums2_i >= 0 else nums1[0] - 1
            if nums1_number >= nums2_number:
                nums1[i] = nums1_number
                nums1_i -= 1
                continue
            
            nums1[i] = nums2_number
            nums2_i -= 1