class Solution {
public:
    int maximumElementAfterDecrementingAndRearranging(vector<int>& arr) {
        std::sort(arr.begin(), arr.end());
        auto length{arr.size()};
        int prev{1};
        for (int i = 1; i < length; ++i) {
            if (arr[i] != prev) {
                ++prev;
            }
        }
        return prev;
    }
};