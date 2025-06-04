class BrowserHistory:

    def __init__(self, homepage: str):
        self._history = [homepage]
        self._cur = 0

    def visit(self, url: str) -> None:
        self._history[self._cur+1:] = (url, )
        self._cur += 1

    def back(self, steps: int) -> str:
        self._cur = max(0, self._cur - steps)
        return self._history[self._cur]

    def forward(self, steps: int) -> str:
        self._cur = min(len(self._history) - 1, self._cur + steps)
        return self._history[self._cur]


# Your BrowserHistory object will be instantiated and called as such:
# obj = BrowserHistory(homepage)
# obj.visit(url)
# param_2 = obj.back(steps)
# param_3 = obj.forward(steps)