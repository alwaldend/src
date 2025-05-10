class TrieNode:
    def __init__(self):
        self.children = {}
        self.is_word = False

class WordDictionary:

    def __init__(self):
        self.root = TrieNode()

    def addWord(self, word: str) -> None:
        cur = self.root
        for char in word:
            cur = cur.children.setdefault(char, TrieNode())
        cur.is_word = True

    def search(self, word: str) -> bool:
        cur, nxt = [self.root], []
        
        for char in word:
            if not cur:
                return False

            if char == ".":
                for node in cur:
                    nxt.extend(node.children.values())
            else:
                for node in cur:
                    if char in node.children:
                        nxt.append(node.children[char]) 
                
            cur.clear()
            cur, nxt = nxt, cur

        return any(node.is_word for node in cur)
            



# Your WordDictionary object will be instantiated and called as such:
# obj = WordDictionary()
# obj.addWord(word)
# param_2 = obj.search(word)