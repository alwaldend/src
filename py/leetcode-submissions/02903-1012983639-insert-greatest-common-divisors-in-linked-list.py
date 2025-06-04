# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def insertGreatestCommonDivisors(self, head: Optional[ListNode]) -> Optional[ListNode]:
        root = head
        
        while head and head.next:
            next = head.next
            new_node = ListNode(math.gcd(head.val, head.next.val), next)
            head.next = new_node
            head = next
        
        return root