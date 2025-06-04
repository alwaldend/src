class Solution:
    def getIntersectionNode(self, headA: ListNode, headB: ListNode) -> Optional[ListNode]:
        if not headA or not headB:
            return None
        
        tail1, tail2 = headA, headB
        while tail1 != tail2:
            tail1 = tail1.next if tail1 else headB
            tail2 = tail2.next if tail2 else headA
        
        return tail2
    