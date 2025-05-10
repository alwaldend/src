class Solution:
    def twoSum(self, numbers: List[int], target: int) -> List[int]:
        numbers_map: Dict[str, int] = {}

        for i, number in enumerate(numbers):
            diff = target - number
            if diff in numbers_map:
                return [numbers_map[diff] + 1, i + 1]
            numbers_map[number] = i
        
        return []
