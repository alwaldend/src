class Solution:
    def isValid(self, s: str) -> bool:
        brackets = []

        symbols = {
            "{": "}",
            "(": ")",
            "[": "]"
        }

        open = symbols.keys()
        closed = symbols.values()

        for bracket in s:
            if bracket in open:
                brackets.append(bracket)
                continue

            if not brackets:
                return False

            last_bracket = brackets[-1]
            if last_bracket in closed: 
                return False

            correct_closing_bracket = symbols[last_bracket]
            if bracket != correct_closing_bracket:
                return False
            
            brackets.pop()

        return not len(brackets)