/* https://leetcode.com/problems/binary-search/

Given an array of integers nums which is sorted in ascending order, and an integer target, write a function to search target in nums. If target exists, then return its index. Otherwise, return -1.

You must write an algorithm with O(log n) runtime complexity.



Example 1:

Input: nums = [-1,0,3,5,9,12], target = 9
Output: 4
Explanation: 9 exists in nums and its index is 4

Example 2:

Input: nums = [-1,0,3,5,9,12], target = 2
Output: -1
Explanation: 2 does not exist in nums so return -1



Constraints:

    1 <= nums.length <= 104
    -104 < nums[i], target < 104
    All the integers in nums are unique.
    nums is sorted in ascending order.

*/
package main

func search(numbers []int, target int) int {
	length := len(numbers)
	switch {
	case numbers[0] > target:
		// smallet number > target = there is no target in the array
		fallthrough
	case numbers[length-1] < target:
		// biggest number < target = there is no target in the array
		return -1
	case numbers[length-1] == target:
		// checking just in case, can save time
		return length - 1
	case numbers[0] == target:
		// checking just in case, can save time
		return 0
	}

	left, right := 0, length-1
	for right >= left {
		// overflow protection
		index := left + (right-left)/2
		number := numbers[index]
		switch {
		case number == target:
			// found the target
			return index
		case number > target:
			// array is in the ascending order, the number is bigger
			// -> the target is to the left -> discard right
			right = index - 1
		case number < target:
			// array is in the ascending order, the number is smaller
			// -> the target is to the right -> discard left
			left = index + 1
		}
	}
	// search space is empty, there is no target
	return -1
}
