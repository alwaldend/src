# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def partition(self, head: Optional[ListNode], x: int) -> Optional[ListNode]:
        before_head, after_head = ListNode(0), ListNode(0)
        before_tail, after_tail = before_head, after_head
        
        while head: 
            if head.val < x:
                before_tail.next, before_tail = head, head
            else:
                after_tail.next, after_tail = head, head
            head = head.next
        
        after_tail.next, before_tail.next = None, after_head.next
        
        return before_head.next