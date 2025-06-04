class Solution:
    def evalRPN(self, tokens: List[str]) -> int:
        stack: list[int] = []
        for token in tokens:
            match token:
                case "+":
                    stack.append(stack.pop() + stack.pop())
                case "-":
                    last, prev = stack.pop(), stack.pop()
                    stack.append(prev - last)
                case "*":
                    stack.append(stack.pop() * stack.pop())
                case "/": 
                    last, prev = stack.pop(), stack.pop()
                    stack.append(int(prev / last))
                case _:
                    stack.append(int(token))
        return stack[0]