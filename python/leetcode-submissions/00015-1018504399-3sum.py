class Solution:
    def threeSum(self, nums: List[int]) -> List[List[int]]:
        result = set()

        #1. Split nums into three lists: negative numbers, positive numbers, and zeros
        negatives, positives, zeros = [], [], []
        for num in nums:
            if num > 0:
                positives.append(num)
            elif num < 0: 
                negatives.append(num)
            else:
                zeros.append(num)

        #2. Create a separate set for negatives and positives for O(1) look-up times
        negatives_count, positives_count, zeros_count = len(negatives), len(positives), len(zeros)
        negatives_set, positives_set = set(negatives), set(positives)

        #3. If there is at least 1 zero in the list, add all cases where -num exists in N and num exists in P
        #   i.e. (-3, 0, 3) = 0
        for num in positives_set if zeros else []:
            negative = -1 * num
            if negative in negatives_set:
                result.add((negative, 0, num))

        #3. If there are at least 3 zeros in the list then also include (0, 0, 0) = 0
        if zeros_count >= 3:
            result.add((0, 0, 0))

        #4. For all pairs of negative numbers (-3, -1), check to see if their complement (4)
        #   exists in the positive number set
        for i in range(negatives_count):
            negative_1 = negatives[i]
            for j in range(i + 1, negatives_count):
                negative_2 = negatives[j]
                target = -1 * (negative_1 + negative_2)
                if target in positives_set:
                    result.add(tuple(sorted([negative_1, negative_2, target])))

        #5. For all pairs of positive numbers (1, 1), check to see if their complement (-2)
        #   exists in the negative number set
        for i in range(positives_count):
            positive_1 = positives[i]
            for j in range(i + 1, positives_count):
                positive_2 = positives[j]
                target = -1 * (positive_1 + positive_2)
                if target in negatives_set:
                    result.add(tuple(sorted([positive_1, positive_2, target])))

        return result