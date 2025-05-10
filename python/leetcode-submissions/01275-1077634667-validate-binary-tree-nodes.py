class Solution:
    def validateBinaryTreeNodes(self, n: int, leftChild: List[int], rightChild: List[int]) -> bool:
        parents = [-1] * n

        for i in range(n):
            left, right = leftChild[i], rightChild[i] 
            left_valid, right_valid = left != -1, right != -1
            if (left_valid and parents[left] != -1) or (
                right_valid and parents[right] != -1
            ):
                return False
            if left_valid:
                parents[left] = i
            if right_valid:
                parents[right] = i

            parent = parents[i]
            if parent != -1 and (parent == left or parent == right):
                return False

        root = None
        for i, node in enumerate(parents):
            if node == -1 and root is not None:
                return False
            if node == -1:
                root = i

        if root is None:
            return False

        to_visit = [root]
        visited = set()

        while to_visit:
            node = to_visit.pop()
            if node in visited:
                return False
            visited.add(node)
            left, right = leftChild[node], rightChild[node]
            if left != -1:
                to_visit.append(left)
            if right != -1:
                to_visit.append(right)

        return len(visited) == n 
        
