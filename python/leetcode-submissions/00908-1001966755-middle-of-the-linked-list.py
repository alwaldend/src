# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def middleNode(self, head: Optional[ListNode]) -> Optional[ListNode]:
        middle = head.next if head else None 
        tail = middle.next if middle else None
        move = True
        while tail and tail.next:
            tail = tail.next
            if move:
                middle = middle.next
            
            move = not move

        return middle if middle else head
