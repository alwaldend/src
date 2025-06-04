class ATM:

    def __init__(self):
        self._banknotes = [0] * 5
        self._values = (20, 50, 100, 200, 500)

    def deposit(self, banknotesCount: List[int]) -> None:
        for i in range(5):
            self._banknotes[i] += banknotesCount[i]

    def withdraw(self, amount: int) -> List[int]:
        withdrawn = [0] * 5
        for i in reversed(range(5)):
            value, notes_left = self._values[i], self._banknotes[i]
            notes_need = min(notes_left, amount // value)
            
            if notes_need == 0:
                continue
            
            amount -= value * notes_need
            withdrawn[i] = notes_need

            if amount == 0:
                for i in range(5):
                    self._banknotes[i] -= withdrawn[i]
                return withdrawn

        return [-1]


# Your ATM object will be instantiated and called as such:
# obj = ATM()
# obj.deposit(banknotesCount)
# param_2 = obj.withdraw(amount)