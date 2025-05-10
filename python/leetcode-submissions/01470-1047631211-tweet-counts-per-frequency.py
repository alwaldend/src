class TweetCounts:

    def __init__(self):
        self._tweets = defaultdict(list)
        self._chunk_ranges = {
            "minute": 60, 
            "hour": 3600, 
            "day": 86400
        }

    def recordTweet(self, tweetName: str, time: int) -> None:
        self._tweets[tweetName].append(time)

    def getTweetCountsPerFrequency(self, freq: str, tweetName: str, startTime: int, endTime: int) -> List[int]:
        chunk_range = self._chunk_ranges[freq]
        chunks = [0] * (1 + (endTime - startTime) // chunk_range)

        for tweet in self._tweets[tweetName]:
            if not startTime <= tweet <= endTime:
                continue
            chunks[(tweet - startTime) // chunk_range] += 1
        
        return chunks




# Your TweetCounts object will be instantiated and called as such:
# obj = TweetCounts()
# obj.recordTweet(tweetName,time)
# param_2 = obj.getTweetCountsPerFrequency(freq,tweetName,startTime,endTime)