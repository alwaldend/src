/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func mergeInBetween(list1 *ListNode, a int, b int, list2 *ListNode) *ListNode {
    if list1 == nil {
        return list2 
    }
    list2Head, list2Tail := list2, list2
    for list2 != nil {
        list2Tail = list2
        list2 = list2.Next
    }
    count := 0
    list1Head := list1
    for list1 != nil {
        next := list1.Next
        if count == a - 1 {
            list1.Next = list2Head
        } else if count == b {
            list2Tail.Next = list1.Next
        }
        list1 = next
        count++
    }
    return list1Head
}