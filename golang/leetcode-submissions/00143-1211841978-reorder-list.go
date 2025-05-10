func reorderList(head *ListNode) {
    if head == nil || head.Next == nil {
        return
    }

    middle := midNode(head)
    newHead := middle.Next
    middle.Next = nil

    newHead = reverseLinkedList(newHead)

    c1 := head
    c2 := newHead
    var f1, f2 *ListNode

    for c1 != nil && c2 != nil {
        // Backup
        f1 = c1.Next
        f2 = c2.Next

        // Linking
        c1.Next = c2
        c2.Next = f1

        // Move
        c1 = f1
        c2 = f2
    }
}

func midNode(head *ListNode) *ListNode {
    slow := head
    fast := head

    for fast.Next != nil && fast.Next.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    return slow
}

func reverseLinkedList(head *ListNode) *ListNode {
    var prev, curr, forw *ListNode = nil, head, nil

    for curr != nil {
        forw = curr.Next
        curr.Next = prev
        prev = curr
        curr = forw
    }
    return prev
}