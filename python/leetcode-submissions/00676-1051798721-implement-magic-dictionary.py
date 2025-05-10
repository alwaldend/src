class MagicDictionary:

    def __init__(self):
        self.dict = []

    def buildDict(self, dictionary: List[str]) -> None:
        self.dict = dictionary

    def search(self, searchWord: str) -> bool:
        target_length = len(searchWord)
        for word in self.dict:
            if len(word) != target_length or word == searchWord:
                continue
            found_diff = False
            for i in range(target_length):
                if word[i] == searchWord[i]:
                    continue
                if found_diff:
                    break

                found_diff = True
            else:
                return True

        return False

# Your MagicDictionary object will be instantiated and called as such:
# obj = MagicDictionary()
# obj.buildDict(dictionary)
# param_2 = obj.search(searchWord)