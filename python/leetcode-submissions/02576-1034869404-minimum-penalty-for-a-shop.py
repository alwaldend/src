class Solution:
    def bestClosingTime(self, customers: str) -> int:
        # Start with closing at hour 0, the penalty equals all 'Y' in closed hours.
        cur_penalty = min_penalty = customers.count("Y")
        earliest_hour = 0
        
        for hour, customer in enumerate(customers):
            # If status in hour i is 'Y', moving it to open hours decrement
            # penalty by 1. Otherwise, moving 'N' to open hours increment
            # penatly by 1.
            cur_penalty += 1 if customer == "N" else -1

            # Update earliest_hour if a smaller penatly is encountered
            if cur_penalty < min_penalty:
                earliest_hour = hour + 1
                min_penalty = cur_penalty
                
        return earliest_hour