/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func middleNode(head *ListNode) *ListNode {
    if head.Next == nil {
        return head
    }
    slow, fast := head, head.Next
    move := true
    for fast.Next != nil {
        if move {
            slow = slow.Next
            move = false
        } else {
            move = true
        }
        fast = fast.Next
    }
    if move {
        slow = slow.Next
    }
    return slow
}