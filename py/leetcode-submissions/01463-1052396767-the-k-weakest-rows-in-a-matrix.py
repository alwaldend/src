class Solution:
    def kWeakestRows(self, mat: List[List[int]], k: int) -> List[int]:
        rows_count, cols_count = len(mat), len(mat[0])
        heap = []

        for i, row in enumerate(mat):
            count = 0
            for j, val in enumerate(row):
                if val == 0:
                    count = j
                    break
            else:
                count = cols_count
        
            heappush(heap, (count, i))
        
        return tuple(item[1] for item in nsmallest(k, heap))

