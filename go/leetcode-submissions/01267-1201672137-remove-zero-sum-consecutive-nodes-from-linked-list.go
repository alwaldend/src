/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func removeZeroSumSublists(head *ListNode) *ListNode {
    dummy := &ListNode{0, head}
    prefixSumToNode := make(map[int]*ListNode)
    prefixSum := 0
    for current := dummy; current != nil; current = current.Next {
        prefixSum += current.Val
        if prev, found := prefixSumToNode[prefixSum]; found {
            toRemove := prev.Next
            p := prefixSum
            if toRemove != nil {
                p += toRemove.Val
            }
            for toRemove != nil && p != prefixSum {
                delete(prefixSumToNode, p)
                toRemove = toRemove.Next
                if toRemove != nil {
                    p += toRemove.Val
                }
            }
            prev.Next = current.Next
        } else {
            prefixSumToNode[prefixSum] = current
        }
    }
    return dummy.Next
}
