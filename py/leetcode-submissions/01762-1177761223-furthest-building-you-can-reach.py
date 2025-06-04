class Solution:
   def furthestBuilding(self, h: List[int], b: int, l: int) -> int:
       p = []
       
       i = 0
       for i in range(len(h) - 1):
           diff = h[i + 1] - h[i]
           
           if diff <= 0:
               continue
           
           b -= diff
           x = heapq.heappush(p, -diff)
           print(x)
           if b < 0:
               b += -heapq.heappop(p)
               l -= 1
               
           if l < 0:
               return i
       return len(h)-1
  

