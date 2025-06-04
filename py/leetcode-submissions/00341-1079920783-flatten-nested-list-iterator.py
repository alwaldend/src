# """
# This is the interface that allows for creating nested lists.
# You should not implement it, or speculate about its implementation
# """
#class NestedInteger:
#    def isInteger(self) -> bool:
#        """
#        @return True if this NestedInteger holds a single integer, rather than a nested list.
#        """
#
#    def getInteger(self) -> int:
#        """
#        @return the single integer that this NestedInteger holds, if it holds a single integer
#        Return None if this NestedInteger holds a nested list
#        """
#
#    def getList(self) -> [NestedInteger]:
#        """
#        @return the nested list that this NestedInteger holds, if it holds a nested list
#        Return None if this NestedInteger holds a single integer
#        """

class NestedIterator:
    def __init__(self, nestedList: [NestedInteger]):
        self.get_next = self.get_next_gen(nestedList)
        self.next_val = next(self.get_next, None)

    def get_next_gen(self, nestedList: List[NestedInteger]) -> Generator[None, None, int]:
        for ni in nestedList:
            if ni.isInteger():
                yield ni.getInteger()
            else:
                yield from self.get_next_gen(ni.getList())

    def next(self) -> int: 
        answer, self.next_val = self.next_val, next(self.get_next, None)
        return answer
    
    def hasNext(self) -> bool:
        return self.next_val is not None

# Your NestedIterator object will be instantiated and called as such:
# i, v = NestedIterator(nestedList), []
# while i.hasNext(): v.append(i.next())