class Solution:
    def maximumRequests(self, n, requests):
        current = [0] * n
        length = len(requests)
        
        def backtrack(req_index: int, count: int) -> int:
            if req_index == length:
                return 0 if any(current) else count

            req_from, req_to = requests[req_index]
            
            current[req_from] -= 1
            current[req_to] += 1
            take = backtrack(req_index + 1, count + 1)

            current[req_from] += 1
            current[req_to] -= 1
            non_take = backtrack(req_index + 1, count)
            
            return max(take, non_take)
        
        return backtrack(0, 0)