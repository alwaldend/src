class SubrectangleQueries {
public:
    std::vector<std::vector<int>> rectangle;
    std::vector<std::array<int, 5>> updates;
    SubrectangleQueries(vector<vector<int>>& rectangle) {
        this->rectangle = rectangle;
    }
    
    void updateSubrectangle(int row1, int col1, int row2, int col2, int newValue) {
        this->updates.push_back({row1, col1, row2, col2, newValue});
    }
    
    int getValue(int row, int col) {
        int length = this->updates.size();
        for (int i = length - 1; i >= 0; --i) {
            auto const [row1, col1, row2, col2, newValue] = this->updates[i];
            if (row < row1 || row > row2 || col < col1 || col > col2) {
                continue;
            }
            return newValue;
        }
        return this->rectangle[row][col];
    }
};

/**
 * Your SubrectangleQueries object will be instantiated and called as such:
 * SubrectangleQueries* obj = new SubrectangleQueries(rectangle);
 * obj->updateSubrectangle(row1,col1,row2,col2,newValue);
 * int param_2 = obj->getValue(row,col);
 */