# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next

class Solution:
    def mergeTwoLists(self, list1: Optional[ListNode], list2: Optional[ListNode]) -> Optional[ListNode]:
        result = ListNode()
        current = result
        while list1 and list2:               
            if list1.val < list2.val:
                current.next = list1
                list1, current = list1.next, list1
                continue

            current.next = list2
            list2, current = list2.next, list2
                
        if list1 or list2:
            current.next = list1 if list1 else list2
            
        return result.next