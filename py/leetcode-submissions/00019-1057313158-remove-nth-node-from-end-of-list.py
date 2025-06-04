# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def removeNthFromEnd(self, head: Optional[ListNode], n: int) -> Optional[ListNode]:
        cur_length, target_length, before_removed = 0, n + 1, head
        tail = head
        while tail and cur_length != target_length:
            tail = tail.next
            cur_length += 1
        if not tail and cur_length != target_length:
            return head.next
        while tail:
            before_removed = before_removed.next
            tail = tail.next
        before_removed.next = before_removed.next.next
        return head