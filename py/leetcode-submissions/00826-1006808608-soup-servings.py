class Solution:
    def soupServings(self, n: int) -> float:
        servings = ceil(n / 25)

        states = defaultdict(dict)
        moves = [[-4, 0], [-3, -1], [-2, -2], [-1, -3]]
    
        @cache
        def calculate(soup_a: int, soup_b: int) -> float:
            if soup_a <= 0 and soup_b <= 0:
                return 0.5
            if soup_a <= 0:
                return 1.0
            if soup_b <= 0:
                return 0.0
            if soup_a in states and soup_b in states[soup_a]:
                return states[soup_a][soup_b]
            
            state = sum(calculate(soup_a + move[0], soup_b + move[1]) for move in moves) / 4.0
            states[soup_a][soup_b] = state

            return state

        max_probability = 1 - 1e-5

        for serving in range(1, servings + 1):
            state = calculate(serving, serving)
            if state > max_probability:
                return 1.0
        
        return calculate(servings, servings)
       

