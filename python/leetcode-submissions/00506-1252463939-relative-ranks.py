import heapq

class Solution:
    def findRelativeRanks(self, score: List[int]) -> List[str]:
        heap = []
        length = len(score)
        scores = [0] * length
        for i in range(length):
            scores[i] = (score[i], i)
        scores.sort(reverse=True)
        ranks = {0: "Gold Medal", 1: "Silver Medal", 2: "Bronze Medal"}
        for i in range(length):
            _, idx = scores[i]
            score[idx] = ranks.get(i, str(i + 1))
        return score