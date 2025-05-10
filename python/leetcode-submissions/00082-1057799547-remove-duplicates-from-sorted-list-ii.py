# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def deleteDuplicates(self, head: Optional[ListNode]) -> Optional[ListNode]:
        tail, before_tail = head, None
        remove_tail = False
        while tail:
            if tail.next and tail.val == tail.next.val:
                tail.next = tail.next.next
                remove_tail = True
                continue
            if not remove_tail:
                before_tail, tail = tail, tail.next
                continue
            remove_tail = False
            if before_tail is None:
                tail = head.next
                head.next, head = None, head.next
            else:
                before_tail.next = tail.next
                tail = tail.next
        return head