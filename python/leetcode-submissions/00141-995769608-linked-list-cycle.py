# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, x):
#         self.val = x
#         self.next = None

class Solution:
    def hasCycle(self, head: Optional[ListNode]) -> bool:
        slow_p, fast_p = head, head.next if head else None
        while slow_p and fast_p:
            if slow_p == fast_p:
                return True
            slow_p = slow_p.next
            fast_p = fast_p.next if fast_p else None
            fast_p = fast_p.next if fast_p else None
        
        return False