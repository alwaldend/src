class Solution:
    def evalRPN(self, tokens: List[str]) -> int:
        stack = []
        operations = {
            "+": lambda first, second: first + second,
            "-": lambda first, second: first - second,
            "*": lambda first, second: first * second,
            "/": lambda first, second: int(first / second)
        }
        for token in tokens:
            if token not in operations:
                stack.append(int(token))
                continue

            second, first = stack.pop(), stack.pop()
            result = operations[token](first, second)
            stack.append(result)
            
        return stack[-1]