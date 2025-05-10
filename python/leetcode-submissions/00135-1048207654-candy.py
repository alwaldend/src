class Solution:
    def candy(self, ratings: List[int]) -> int:
        child_count = len(ratings)
        candies = [1] * child_count 

        for i in range(1, child_count):
            if ratings[i] > ratings[i-1]:
                candies[i] = candies[i-1] + 1

        for i in reversed(range(child_count - 1)):
            if ratings[i] > ratings[i+1]:
                candies[i] = max(candies[i], candies[i+1] + 1)
        
        return sum(candies)