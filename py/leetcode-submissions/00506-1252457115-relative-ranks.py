import heapq

class Solution:
    def findRelativeRanks(self, score: List[int]) -> List[str]:
        heap = []
        for i, athlete_score in enumerate(score):
            heapq.heappush(heap, (-athlete_score, i))
        i = 0
        ranks = {0: "Gold Medal", 1: "Silver Medal", 2: "Bronze Medal"}
        while heap:
            _, athlete = heapq.heappop(heap)
            rank = ranks.get(i, str(i + 1))
            score[athlete] = rank
            i += 1
        return score