# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def deleteDuplicates(self, head: Optional[ListNode]) -> Optional[ListNode]:
        result = head
        
        while head and head.next:
            if head.next.val == head.val:
                head.next = head.next.next
                continue
            
            head = head.next
        
        return result