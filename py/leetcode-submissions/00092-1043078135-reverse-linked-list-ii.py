# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def reverseBetween(self, head: Optional[ListNode], left: int, right: int) -> Optional[ListNode]:
        cur_node = 1
        left_head, left_tail = head if left > 1 else None, None

        while cur_node < left:
            cur_node += 1
            head, left_tail = head.next, head

        cur_node += 1
        mid_head, mid_tail, head = head, head, head.next
        mid_head.next = None

        while cur_node <= right:
            cur_node += 1
            next_node = head.next
            mid_head, head.next = head, mid_head
            head = next_node
    
        mid_tail.next = head
        if left_head:
            left_tail.next = mid_head
        else:
            left_head = mid_head

        return left_head