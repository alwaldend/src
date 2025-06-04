class Solution:
    def reorganizeString(self, s: str) -> str:
        result = []
        # Min heap ordered by character counts, so we will use
        # negative values for count
        priority_queue = [(-count, char) for char, count in Counter(s).items()]
        heapify(priority_queue)

        while priority_queue:
            count_first, char_first = heappop(priority_queue)
            if not result or char_first != result[-1]:
                result.append(char_first)
                if count_first != -1: 
                    heappush(priority_queue, (count_first + 1, char_first))
                continue
            
            if not priority_queue: 
                return ""
            
            count_second, char_second = heappop(priority_queue)
            result.append(char_second)
            if count_second != -1:
                heappush(priority_queue, (count_second + 1, char_second))
            heappush(priority_queue, (count_first, char_first))

        return "".join(result)