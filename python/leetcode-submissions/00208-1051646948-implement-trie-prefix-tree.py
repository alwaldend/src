class Node:
    def __init__(self):
        self.ch, self.is_word = [None] * 128, False

class Trie:

    def __init__(self):
        self.head = Node()

    def insert(self, word: str) -> None:
        cur_node = self.head
        for i in range(len(word)):
            char = ord(word[i])
            if cur_node.ch[char] is None:
                cur_node.ch[char] = Node()
            cur_node = cur_node.ch[char]
        cur_node.is_word = True

    def search(self, word: str) -> bool:
        cur_node = self.head
        for i in range(len(word)):
            char = ord(word[i])
            if cur_node.ch[char] is None:
                return False
            cur_node = cur_node.ch[char]
        
        return cur_node.is_word

    def startsWith(self, prefix: str) -> bool:
        cur_node = self.head
        for i in range(len(prefix)):
            char = ord(prefix[i])
            if cur_node.ch[char] is None:
                return False
            cur_node = cur_node.ch[char]
        
        return True


# Your Trie object will be instantiated and called as such:
# obj = Trie()
# obj.insert(word)
# param_2 = obj.search(word)
# param_3 = obj.startsWith(prefix)