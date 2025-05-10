/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func reverseList(head *ListNode) *ListNode {
    if head == nil {
        return nil
    }
    cur, next := head, head.Next
    head.Next = nil
    for next != nil {
        cur, next, next.Next = next, next.Next, cur
    }
    return cur
}