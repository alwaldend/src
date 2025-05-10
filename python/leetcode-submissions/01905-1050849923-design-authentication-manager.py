class AuthenticationManager:

    def __init__(self, timeToLive: int):
        self.ttl = timeToLive
        self.tokens = {}

    def generate(self, tokenId: str, currentTime: int) -> None:
        self.tokens[tokenId] = currentTime + self.ttl

    def renew(self, tokenId: str, currentTime: int) -> None:
        if currentTime < self.tokens.get(tokenId, currentTime):
            self.generate(tokenId, currentTime)
        else:
            self.tokens.pop(tokenId, None)

    def countUnexpiredTokens(self, currentTime: int) -> int:
        count = 0
        remove_tokens = []
        for token, expir_time in self.tokens.items():
            if currentTime < expir_time:
                count += 1
            else:
                remove_tokens.append(token)
        
        for token in remove_tokens:
            self.tokens.pop(token)
        
        return count


# Your AuthenticationManager object will be instantiated and called as such:
# obj = AuthenticationManager(timeToLive)
# obj.generate(tokenId,currentTime)
# obj.renew(tokenId,currentTime)
# param_3 = obj.countUnexpiredTokens(currentTime)