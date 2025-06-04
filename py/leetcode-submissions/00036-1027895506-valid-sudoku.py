class Solution:
    def isValidSudoku(self, board: List[List[str]]) -> bool:
        row_count, column_count = 9, 9

        

        row_counters = [defaultdict(bool) for _ in range(row_count)]
        column_counters = [defaultdict(bool) for _ in range(column_count)]
        subbox_counters = [[defaultdict(bool) for _ in range(column_count//3)] 
                            for _ in range(row_count//3)]

        for row in range(row_count):
            for column in range(column_count):
                char = board[row][column]
                if char == ".":
                    continue

                counters = (
                    row_counters[row], column_counters[column], 
                    subbox_counters[row//3][column//3]
                )
                for counter in counters:
                    if counter[char]:
                        return False

                    counter[char] = True
        
        return True

        