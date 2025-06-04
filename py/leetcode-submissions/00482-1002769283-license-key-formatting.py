class Solution:
    # "5F3Z-2e-9-w" -> "5F3Z2E9W" -> "5F3Z2E9W"
    def licenseKeyFormatting(self, s: str, k: int) -> str:
        letters = s.replace("-", "").upper()
        result = []
        end = len(letters)
        while end > 0:
            start = end - k
            if start < 0:
                start = 0
            result.append(letters[start:end])
            end -= k

        return "-".join(reversed(result))