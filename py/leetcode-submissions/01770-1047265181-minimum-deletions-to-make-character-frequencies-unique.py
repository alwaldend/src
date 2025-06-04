class Solution:
    def minDeletions(self, s: str) -> int:
        length = len(s)
        char_count = Counter(s)
        counts = set()
        deleted_chars = 0

        for char, count in char_count.items():
            while count != 0 and count in counts:
                count -= 1
                deleted_chars += 1
            if count != 0:
                counts.add(count)
            
        return deleted_chars
