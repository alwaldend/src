# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def reverse(self, head: ListNode) -> ListNode:
        previous = None
        while head:
            next, head.next = head.next, previous
            previous, head = head, next
        return previous

    def addTwoNumbers(self, l1: Optional[ListNode], l2: Optional[ListNode]) -> Optional[ListNode]:
        l1, l2 = self.reverse(l1), self.reverse(l2)
        head, tail, carry = l1, l1, 0

        while l1 or l2 or carry:
            val_1 = l1.val if l1 else 0
            val_2 = l2.val if l2 else 0
            sum = val_1 + val_2 + carry
            if sum > 9:
                carry = 1
                sum %= 10
            else:
                carry = 0

            tail.val = sum
            l1, l2 = l1.next if l1 else None, l2.next if l2 else None
            if not tail.next and l2:
                tail.next = l2
                l1 = None

            if not tail.next and carry:
                tail.next = ListNode(carry)
                carry = 0
                l2 = None

            if tail.next:
                tail = tail.next
            
        return self.reverse(head)