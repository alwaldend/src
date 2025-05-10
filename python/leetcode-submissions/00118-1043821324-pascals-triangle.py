class Solution:
    def generate(self, numRows: int) -> List[List[int]]:
        answer = [[1]]
        numRows -= 1

        while numRows > 0:
            cur, prev = [1], answer[-1]
            for i in range(len(prev) - 1):
                cur.append(prev[i] + prev[i+1])
            cur.append(1)
            answer.append(cur)
            numRows -= 1
        
        return answer