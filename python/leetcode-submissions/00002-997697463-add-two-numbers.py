# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    # create the result pointer that points towards one of the lists
    # create the result root pointer that will point towards the root of the result list
    # 1. start infinite loop
    # 2. if result is false - break, we reached the end
    # 3. get values from non-empty pointers
    # 4. add values together, store the carry in a variable
    # 5. store the value in the result pointer
    # 6. move the result pointer to the next node, if there is no next node, use a node from
    #    another list
    # 7. move list pointers
    # 8. if we have a carry left, add a node to the result list
    def addTwoNumbers(self, l1: Optional[ListNode], 
                            l2: Optional[ListNode]) -> Optional[ListNode]:
        carry = 0
        result, result_root = l1, l1

        while True:
            if not l1 and not l2:
                break
            
            number_1 = l1.val if l1 else 0
            number_2 = l2.val if l2 else 0
            result_val = number_1 + number_2 + carry
            if result_val > 9:
                carry, result_val = 1, result_val - 10
            else:
                carry = 0
            
            result.val = result_val
            if not result.next and l2:
                l1 = None
                result.next = l2.next

            if result.next:
                result = result.next
            if l1:
                l1 = l1.next
            if l2:
                l2 = l2.next

        if carry:
            result.next = ListNode(carry)

        return result_root