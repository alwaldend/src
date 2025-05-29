
func sortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	if head.Next.Next == nil {
		a := head
		b := head.Next
		if a.Val <= b.Val {
			a.Next = b
			b.Next = nil
			return a
		} else {
			b.Next = a
			a.Next = nil
			return b
		}
	}
	slow := head
	fast := head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	nextList := slow.Next
	slow.Next = nil
	list1 := sortList(head)
	list2 := sortList(nextList)
	dummy := &ListNode{Val: -1, Next: nil}
	it := dummy
	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			it.Next = list1
			list1 = list1.Next
		} else {
			it.Next = list2
			list2 = list2.Next
		}
		it = it.Next
	}
	if list1 != nil {
		it.Next = list1
	} else {
		it.Next = list2
	}
	return dummy.Next
}
