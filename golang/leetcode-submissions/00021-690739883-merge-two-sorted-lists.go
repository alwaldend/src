func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	// ensure both lists are valid
	switch {
	case list1 == nil && list2 == nil:
		return nil
	case list1 == nil && list2 != nil:
		return list2
	case list1 != nil && list2 == nil:
		return list1
	}
	var root *ListNode
	if list1.Val < list2.Val {
		root = &ListNode{list1.Val, nil}
		list1 = list1.Next
	} else {
		root = &ListNode{list2.Val, nil}
		list2 = list2.Next
	}
	current := &root
	for {
		switch {
		case list1 == nil && list2 == nil:
			return root
		case list1 != nil && list2 != nil && list1.Val <= list2.Val:
			fallthrough
		case list1 != nil && list2 == nil:
			fmt.Println("1", list1.Val)
			updateResult(&current, &list1)
		case list1 != nil && list2 != nil && list1.Val > list2.Val:
			fallthrough
		case list1 == nil && list2 != nil:
			fmt.Println("2", list2.Val)
			updateResult(&current, &list2)
		}
	}
}
func updateResult(current ***ListNode, node **ListNode) {
	// modify the current result node
	(**current).Next = &ListNode{(*node).Val, nil}
	// move the current pointer
	*current = &(**current).Next
	// there is no next node -> nill it
	// there is next -> move it
	if (*node).Next == nil {
		*node = nil
	} else {
		*node = (*node).Next
	}
}
