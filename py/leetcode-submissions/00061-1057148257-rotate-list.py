# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def rotateRight(self, head: Optional[ListNode], k: int) -> Optional[ListNode]:
        if k == 0 or not head:
            return head
        temp_head, tail, length = head, head, 0
        while temp_head:
            length += 1
            tail, temp_head = temp_head, temp_head.next
        rotation = k % length
        if rotation == 0:
            return head
        new_head_idx = length - rotation
        new_tail_idx = new_head_idx - 1 
        i, new_head, new_tail = 0, None, None
        temp_head = head
        while temp_head:
            if i == new_head_idx:
                new_head = temp_head
                break
            if i == new_tail_idx:
                new_tail = temp_head
            i += 1
            temp_head = temp_head.next
        tail.next = head
        new_tail.next = None
        return new_head

            
