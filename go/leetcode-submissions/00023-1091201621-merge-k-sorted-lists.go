/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func mergeKLists(lists []*ListNode) *ListNode {
    n:=len(lists)
    if n==0{
        return nil
    }
    curr:=lists[0]
    if n==1{
        return curr
    }
    for i:=1;i<n;i++{
        curr=mergeList(curr,lists[i])
    }
    return curr
}

func mergeList(l1,l2 *ListNode) *ListNode {
    head:=&ListNode{}
    curr:=head
    for l1!=nil && l2!=nil{
        if l1.Val<l2.Val{
            curr.Next=l1
            l1=l1.Next
            curr=curr.Next
        }else{
            curr.Next=l2
            l2=l2.Next
            curr=curr.Next
        }
    }
    if l1 != nil {
        curr.Next = l1
    } else if l2 != nil {
        curr.Next = l2
    }
    return head.Next
}