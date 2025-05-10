func findPeakElement(nums []int) int {
    left := 0
    right := len(nums) - 1
    
    for left < right {
        mid := left + (right - left) / 2
        
        if nums[mid] > nums[mid+1] {
            // The peak is in the left half
            right = mid
        } else {
            // The peak is in the right half
            left = mid + 1
        }
    }
    
    return left
}