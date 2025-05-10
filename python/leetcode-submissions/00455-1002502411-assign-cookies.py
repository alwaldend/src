class Solution:
    def findContentChildren(self, g: List[int], s: List[int]) -> int:
        greed = 0
        g.sort(reverse=True)
        s.sort(reverse=True)
        children_count = len(g)
        count = 0
        for cookie_size in s:
            while greed < children_count and g[greed] > cookie_size:
                greed += 1
            
            if greed >= children_count:
                break

            count += 1
            greed += 1
            
        return count