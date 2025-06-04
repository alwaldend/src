class Solution:
    def splitListToParts(self, head: Optional[ListNode], k: int) -> List[Optional[ListNode]]:        
        nodes_count = 0
        tail = head
        while tail:
            tail = tail.next
            nodes_count += 1
        
        quotient, remainder = divmod(nodes_count, k)
        answer = []
        for i in range(k):
            size = quotient + 1 if (remainder := remainder - 1) >= 0 else quotient

            answer.append(head)
            last = None
            while size > 0:
                last = head
                head = head.next
                size -= 1
            if last: 
                last.next = None
                
        return answer