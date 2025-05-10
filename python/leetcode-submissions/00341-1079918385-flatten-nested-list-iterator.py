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
        self.list = []
        self.flatten(self.list, nestedList)
        self.cur = 0
        self.length = len(self.list)

    def flatten(self, flatList: List[int], nestedList: List[NestedInteger]):
        for ni in nestedList:
            if ni.isInteger():
                flatList.append(ni.getInteger())
            else:
                self.flatten(flatList, ni.getList())

    def next(self) -> int: 
        answer = self.list[self.cur]
        self.cur += 1
        return answer
    
    def hasNext(self) -> bool:
        return self.cur < self.length

# Your NestedIterator object will be instantiated and called as such:
# i, v = NestedIterator(nestedList), []
# while i.hasNext(): v.append(i.next())