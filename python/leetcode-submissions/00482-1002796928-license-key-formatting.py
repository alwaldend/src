class Solution:
    def licenseKeyFormatting(self, s: str, k: int) -> str:
        result = []
        count = 0
        for letter in reversed(s):
            if letter == "-":
                continue

            if count == k:
                result.append("-")
                count = 0
            
            count += 1
            result.append(letter.upper())

        return "".join(reversed(result))