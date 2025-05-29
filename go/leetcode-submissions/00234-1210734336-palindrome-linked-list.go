/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func isPalindrome(head *ListNode) bool {
    nodes := []int{}
    for head != nil {
        nodes = append(nodes, head.Val)
        head = head.Next
    }
    length := len(nodes)
    for i := range(length / 2) {
        start, end := nodes[i], nodes[length-i-1]
        if start != end {
            return false
        }
    }
    return true
}