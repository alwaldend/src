class Solution:
    # [3,0,6,1,5]
    # [0,1,3,5,6]
    def hIndex(self, citations: List[int]) -> int:
        citations = sorted(citations)
        length = len(citations)

        h = 0
        for i in reversed(range(length)):
            citations_count = citations[i]
            published_count = length - i

            if citations_count == 0 or published_count < h:
                break

            if published_count <= citations_count:
                h = published_count
            
        return h

