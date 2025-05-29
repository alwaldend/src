/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    res := &ListNode{0, head}

    lead := res
    for i := 0; i <= n; i++ {
        lead = lead.Next
    }

    cur := res
    for lead != nil {
        cur = cur.Next
        lead = lead.Next
    }

    cur.Next = cur.Next.Next
    return res.Next
}