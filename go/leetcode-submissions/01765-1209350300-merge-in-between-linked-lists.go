/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func mergeInBetween(list1 *ListNode, a int, b int, list2 *ListNode) *ListNode {
    l := list1
    for i := 1; i <= a - 1; i++ {
        l = l.Next
    }

    prev := l
    r := l
    for i := a; i <= b + 1; i++ {
        r = r.Next
        prev.Next = nil
        prev = r
    }

    list2Tail := list2
    for list2Tail.Next != nil {
        list2Tail = list2Tail.Next
    }

    l.Next = list2
    list2Tail.Next = r

    return list1
}