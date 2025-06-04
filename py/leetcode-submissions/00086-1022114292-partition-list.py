# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def partition(self, head: Optional[ListNode], x: int) -> Optional[ListNode]:
        first_left, last_left, first_right, last_right = None, None, None, None

        while head:
            is_left = head.val < x
            is_right = not is_left
            next_node = head.next
            head.next = first_right if is_left else None
            
            if is_left and last_left:
                last_left.next = head
                last_left = head
            elif is_left and not first_left:
                first_left, last_left = head, head

            if is_right and first_right:
                last_right.next = head
                last_right = head
            elif is_right and not first_right:
                first_right, last_right = head, head
                if last_left:
                    last_left.next = head

            head = next_node

        return first_left if first_left else first_right