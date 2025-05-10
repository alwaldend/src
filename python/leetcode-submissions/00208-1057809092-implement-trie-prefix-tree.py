class Trie:

    def __init__(self):
        self.root = {}

    def insert(self, word: str) -> None:
        cur = self.root
        for char in word:
            if char not in cur:
                cur[char] = {}
            cur = cur[char]
        
        cur["_is_word"] = None

    def search(self, word: str) -> bool:
        cur = self.root
        for char in word:
            if char not in cur:
                return False
            cur = cur[char]
        
        return "_is_word" in cur

    def startsWith(self, prefix: str) -> bool:
        cur = self.root
        for char in prefix:
            if char not in cur:
                return False
            cur = cur[char]
            
        return True


# Your Trie object will be instantiated and called as such:
# obj = Trie()
# obj.insert(word)
# param_2 = obj.search(word)
# param_3 = obj.startsWith(prefix)